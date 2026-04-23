package notify

import (
	"fmt"
	"sync"
	"time"
)

type authFailCounter struct {
	mu          sync.Mutex
	count       int
	windowStart time.Time
}

// authFailures tracks per-IP authentication failure bursts.
var authFailures sync.Map // key: IP string, value: *authFailCounter

// TrackAuthFailure increments the failure counter for the given IP.
// When 5 failures accumulate within 60 seconds, a notification is sent and
// the counter resets. Non-blocking: counter update runs in a background goroutine.
func TrackAuthFailure(n Notifier, ip string) {
	go func() {
		v, _ := authFailures.LoadOrStore(ip, &authFailCounter{})
		c := v.(*authFailCounter)
		c.mu.Lock()
		now := time.Now()
		if now.Sub(c.windowStart) > 60*time.Second {
			c.count = 0
			c.windowStart = now
		}
		c.count++
		count := c.count
		c.mu.Unlock()

		if count >= 5 {
			authFailures.Delete(ip)
			FireNotify(n,
				"git-cone: auth failure burst",
				fmt.Sprintf("auth: %d failed attempts from %s in 60s", count, ip),
				7,
			)
		}
	}()
}
