package main

import (
	"fmt"
	"time"
)

type MatterIntake struct {
	Email, MatterID, SignedDocumentURL string
	Deadline                           time.Time
}

func NeedsDeadlineFollowUp(now, deadline time.Time) bool { return !now.Before(deadline) }

func SendMatterSequence(c *InfraiClient, m MatterIntake, now time.Time) ([]string, error) {
	ids := make([]string, 0, 3)
	id, err := c.SendEmail(m.Email, "Verify your legal matter email", fmt.Sprintf("<p>Verify matter %s: <a href=\"%s\">open link</a></p>", m.MatterID, m.SignedDocumentURL))
	if err != nil {
		return nil, err
	}
	ids = append(ids, id)
	if _, err = c.GetEmail(id); err != nil {
		return nil, err
	}
	id, err = c.SendEmail(m.Email, "Signed document ready", fmt.Sprintf("<p>Your signed document for matter %s is ready.</p>", m.MatterID))
	if err != nil {
		return nil, err
	}
	ids = append(ids, id)
	if NeedsDeadlineFollowUp(now, m.Deadline) {
		id, err = c.SendEmail(m.Email, "Matter deadline follow-up", fmt.Sprintf("<p>Please review matter %s before its deadline.</p>", m.MatterID))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
