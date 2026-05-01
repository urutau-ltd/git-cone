package cmd

import (
	"context"
	"net"
	"testing"

	charmssh "github.com/charmbracelet/ssh"
	"github.com/urutau-ltd/git-cone/pkg/sshpolicy"
	gossh "golang.org/x/crypto/ssh"
)

type testAlgorithmsConn struct {
	algs gossh.NegotiatedAlgorithms
}

func (c testAlgorithmsConn) User() string                           { return "" }
func (c testAlgorithmsConn) SessionID() []byte                      { return nil }
func (c testAlgorithmsConn) ClientVersion() []byte                  { return nil }
func (c testAlgorithmsConn) ServerVersion() []byte                  { return nil }
func (c testAlgorithmsConn) RemoteAddr() net.Addr                   { return &net.TCPAddr{} }
func (c testAlgorithmsConn) LocalAddr() net.Addr                    { return &net.TCPAddr{} }
func (c testAlgorithmsConn) Algorithms() gossh.NegotiatedAlgorithms { return c.algs }

func (c testAlgorithmsConn) SendRequest(string, bool, []byte) (bool, []byte, error) {
	return false, nil, nil
}

func (c testAlgorithmsConn) OpenChannel(string, []byte) (gossh.Channel, <-chan *gossh.Request, error) {
	return nil, nil, nil
}

func (c testAlgorithmsConn) Close() error { return nil }
func (c testAlgorithmsConn) Wait() error  { return nil }

func TestNegotiatedAlgorithmsFromContext(t *testing.T) {
	t.Parallel()

	want := gossh.NegotiatedAlgorithms{
		KeyExchange: gossh.KeyExchangeMLKEM768X25519,
		HostKey:     gossh.KeyAlgoED25519,
		Read:        gossh.DirectionAlgorithms{Cipher: gossh.CipherChaCha20Poly1305},
	}
	ctx := context.WithValue(context.Background(), charmssh.ContextKeyConn, testAlgorithmsConn{algs: want})

	got, ok := negotiatedAlgorithmsFromContext(ctx)
	if !ok {
		t.Fatal("expected negotiated algorithms in context")
	}
	if got.KeyExchange != want.KeyExchange {
		t.Fatalf("unexpected kex: %q", got.KeyExchange)
	}
	if got.HostKey != want.HostKey {
		t.Fatalf("unexpected host key: %q", got.HostKey)
	}
	if got.Read.Cipher != want.Read.Cipher {
		t.Fatalf("unexpected cipher: %q", got.Read.Cipher)
	}
}

func TestNegotiatedAlgorithmsFromWrappedServerConn(t *testing.T) {
	t.Parallel()

	want := gossh.NegotiatedAlgorithms{
		KeyExchange: gossh.KeyExchangeMLKEM768X25519,
		HostKey:     gossh.KeyAlgoED25519,
		Read:        gossh.DirectionAlgorithms{Cipher: gossh.CipherChaCha20Poly1305},
	}
	ctx := context.WithValue(context.Background(), charmssh.ContextKeyConn, &gossh.ServerConn{
		Conn: testAlgorithmsConn{algs: want},
	})

	got, ok := negotiatedAlgorithmsFromContext(ctx)
	if !ok {
		t.Fatal("expected negotiated algorithms in wrapped server conn")
	}
	if got.KeyExchange != want.KeyExchange {
		t.Fatalf("unexpected kex: %q", got.KeyExchange)
	}
	if got.HostKey != want.HostKey {
		t.Fatalf("unexpected host key: %q", got.HostKey)
	}
	if got.Read.Cipher != want.Read.Cipher {
		t.Fatalf("unexpected cipher: %q", got.Read.Cipher)
	}
}

func TestIsPostQuantumKEX(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kex  string
		want bool
	}{
		{name: "mlkem", kex: gossh.KeyExchangeMLKEM768X25519, want: true},
		{name: "sntrup", kex: "sntrup761x25519-sha512", want: true},
		{name: "curve25519", kex: gossh.KeyExchangeCurve25519, want: false},
		{name: "empty", kex: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sshpolicy.IsPostQuantumKEX(tt.kex); got != tt.want {
				t.Fatalf("IsPostQuantumKEX(%q) = %t, want %t", tt.kex, got, tt.want)
			}
		})
	}
}

func TestAuthMethod(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), charmssh.ContextKeyPermissions, &charmssh.Permissions{
		Permissions: &gossh.Permissions{Extensions: map[string]string{"pubkey-fp": ""}},
	})
	if got := authMethod(ctx, nil); got != "keyless" {
		t.Fatalf("authMethod() = %q, want keyless", got)
	}
}
