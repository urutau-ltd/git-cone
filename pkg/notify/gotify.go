package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Gotify is a Gotify notification client.
type Gotify struct {
	URL      string
	Token    string
	Priority int
	client   *http.Client
}

// NewGotify creates a new Gotify notifier.
func NewGotify(url, token string, priority int) *Gotify {
	return &Gotify{
		URL:      url,
		Token:    token,
		Priority: priority,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// Send implements Notifier.
func (g *Gotify) Send(ctx context.Context, title, message string, priority int) error {
	payload, err := json.Marshal(map[string]any{
		"title":    title,
		"message":  message,
		"priority": priority,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/message?token=%s", g.URL, g.Token),
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gotify: HTTP %d", resp.StatusCode)
	}
	return nil
}
