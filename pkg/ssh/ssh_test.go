package ssh

import (
	"testing"

	"charm.land/log/v2"
	gossh "golang.org/x/crypto/ssh"
)

func TestHardenedServerConfigPrefersPostQuantumKEX(t *testing.T) {
	t.Parallel()

	sc := hardenedServerConfig(log.Default())(nil)
	if sc == nil {
		t.Fatal("expected non-nil server config")
	}

	if len(sc.KeyExchanges) < 4 {
		t.Fatalf("expected hardened KEX list, got %v", sc.KeyExchanges)
	}

	if got := sc.KeyExchanges[0]; got != gossh.KeyExchangeMLKEM768X25519 {
		t.Fatalf("unexpected first KEX: %q", got)
	}

	if got := sc.KeyExchanges[1]; got != gossh.KeyExchangeCurve25519 {
		t.Fatalf("unexpected second KEX: %q", got)
	}
}
