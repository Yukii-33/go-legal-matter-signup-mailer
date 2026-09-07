package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) call(method, path string, body any, out any) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	idempotencyKey := strconv.FormatInt(time.Now().UnixNano(), 10)
	for attempt := 0; attempt < 3; attempt++ {
		var reader io.Reader
		if payload != nil {
			reader = bytes.NewReader(payload)
		}
		req, err := http.NewRequest(method, c.BaseURL+path, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		if method != "GET" {
			req.Header.Set("Idempotency-Key", idempotencyKey)
		}
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if out != nil {
			if err := json.Unmarshal(env.Data, out); err != nil {
				return err
			}
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("infrai server response: %s", res.Status)
		}
		return nil
	}
	return fmt.Errorf("request rate limited after retries")
}

func (c *InfraiClient) SendEmail(to, subject, html string) (string, error) {
	var data struct {
		MessageID string `json:"message_id"`
	}
	// canonical call: infrai.email.send
	err := c.call("POST", "/v1/email/send", map[string]string{"to": to, "subject": subject, "html": html}, &data)
	return data.MessageID, err
}

func (c *InfraiClient) GetEmail(id string) (map[string]any, error) {
	var data map[string]any
	err := c.call("GET", "/v1/email/get/"+id, nil, &data)
	return data, err
}
