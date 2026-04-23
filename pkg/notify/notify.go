package notify

import (
	"context"
	"time"
)

// Notifier sends a notification.
type Notifier interface {
	Send(ctx context.Context, title, message string, priority int) error
}

// Noop is a no-op notifier used when notifications are disabled.
type Noop struct{}

func (Noop) Send(_ context.Context, _, _ string, _ int) error { return nil }

// FireNotify sends a notification in a background goroutine (fire-and-forget).
func FireNotify(n Notifier, title, message string, priority int) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = n.Send(ctx, title, message, priority)
	}()
}
