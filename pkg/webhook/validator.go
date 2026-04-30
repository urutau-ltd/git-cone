package webhook

import (
	"context"
	"errors"
	"net"

	"github.com/urutau-ltd/git-cone/pkg/ssrf"
)

var (
	// ErrInvalidScheme is returned when the webhook URL scheme is not http or https.
	ErrInvalidScheme = errors.New("webhook URL must use http or https scheme")
	// ErrPrivateIP is returned when the webhook URL resolves to a private IP address.
	ErrPrivateIP = errors.New("webhook URL cannot resolve to private or internal IP addresses")
	// ErrInvalidURL is returned when the webhook URL is invalid.
	ErrInvalidURL = errors.New("invalid webhook URL")
)

// ValidateWebhookURL validates that a webhook URL is safe to use.
// It checks:
// - URL is properly formatted
// - Scheme is http or https
// - Hostname does not resolve to private/internal IP addresses
// - Hostname is not localhost or similar.
func ValidateWebhookURL(rawURL string) error {
	return mapSSRFFromShared(ssrf.ValidateURL(context.Background(), rawURL))
}

// ValidateIPBeforeDial validates an IP address before establishing a connection.
// This is used to prevent DNS rebinding attacks.
func ValidateIPBeforeDial(ip net.IP) error {
	return mapSSRFFromShared(ssrf.ValidateIPBeforeDial(ip))
}

func mapSSRFFromShared(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ssrf.ErrPrivateIP):
		return ErrPrivateIP
	case errors.Is(err, ssrf.ErrInvalidScheme):
		return ErrInvalidScheme
	case errors.Is(err, ssrf.ErrInvalidURL):
		return ErrInvalidURL
	default:
		return err
	}
}
