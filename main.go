package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func main() {
	c, err := NewInfraiClient()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var m MatterIntake
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		ids, err := SendMatterSequence(c, m, time.Now())
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"message_ids": ids})
	})
	_ = http.ListenAndServe(":8080", nil)
}
