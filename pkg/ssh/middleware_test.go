package ssh

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/ssh"
	"github.com/matryer/is"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/db/migrate"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/store"
	"github.com/urutau-ltd/git-cone/pkg/store/database"
	gossh "golang.org/x/crypto/ssh"
	_ "modernc.org/sqlite"
)

func TestAuthenticatedUserForPublicKey(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	ctx, cfg, be, _, _, keys := setupAuthTest(t)

	t.Run("known user key resolves user", func(t *testing.T) {
		is := is.New(t)
		user, err := authenticatedUserForPublicKey(ctx, be, cfg, keys.user)
		is.NoErr(err)
		is.True(user != nil)
		is.Equal(user.Username(), "testuser")
	})

	t.Run("bootstrap admin key is accepted without user", func(t *testing.T) {
		is := is.New(t)
		user, err := authenticatedUserForPublicKey(ctx, be, cfg, keys.bootstrapAdmin)
		is.NoErr(err)
		is.True(user == nil)
	})

	t.Run("unknown key is rejected", func(t *testing.T) {
		is := is.New(t)
		user, err := authenticatedUserForPublicKey(ctx, be, cfg, keys.unknown)
		is.True(err == proto.ErrUserNotFound)
		is.True(user == nil)
	})
}

func TestPublicKeyHandlerRejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	ctx, cfg, be, _, _, keys := setupAuthTest(t)
	srv := &SSHServer{cfg: cfg, be: be}

	unknownCtx := newMockSSHContext(ctx)
	is.True(!srv.PublicKeyHandler(unknownCtx, keys.unknown))
	is.True(unknownCtx.permissions.Extensions["pubkey-fp"] == "")

	userCtx := newMockSSHContext(ctx)
	is.True(srv.PublicKeyHandler(userCtx, keys.user))
	is.Equal(userCtx.permissions.Extensions["pubkey-fp"], gossh.FingerprintSHA256(keys.user))

	adminCtx := newMockSSHContext(ctx)
	is.True(srv.PublicKeyHandler(adminCtx, keys.bootstrapAdmin))
	is.Equal(adminCtx.permissions.Extensions["pubkey-fp"], gossh.FingerprintSHA256(keys.bootstrapAdmin))
}

type authTestKeys struct {
	user           gossh.PublicKey
	unknown        gossh.PublicKey
	bootstrapAdmin gossh.PublicKey
}

func setupAuthTest(tb testing.TB) (context.Context, *config.Config, *backend.Backend, *db.DB, store.Store, authTestKeys) {
	tb.Helper()

	is := is.New(tb)
	dp := tb.TempDir()
	cfg := config.DefaultConfig()
	cfg.DataPath = dp
	cfg.DB.Driver = "sqlite"
	cfg.DB.DataSource = filepath.Join(dp, "test.db")

	userPair, userKey := mustGenerateKey(tb, filepath.Join(dp, "user"))
	_ = userPair
	unknownPair, unknownKey := mustGenerateKey(tb, filepath.Join(dp, "unknown"))
	_ = unknownPair
	adminPair, adminKey := mustGenerateKey(tb, filepath.Join(dp, "bootstrap-admin"))
	cfg.InitialAdminKeys = []string{adminPair.AuthorizedKey()}

	ctx := config.WithContext(context.Background(), cfg)
	dbx, err := db.Open(ctx, cfg.DB.Driver, cfg.DB.DataSource)
	is.NoErr(err)
	tb.Cleanup(func() { _ = dbx.Close() })

	is.NoErr(migrate.Migrate(ctx, dbx))
	dbstore := database.New(ctx, dbx)
	ctx = store.WithContext(ctx, dbstore)
	be := backend.New(ctx, cfg, dbx, dbstore)
	ctx = backend.WithContext(ctx, be)

	_, err = be.CreateUser(ctx, "testuser", proto.UserOptions{
		PublicKeys: []gossh.PublicKey{userKey},
	})
	is.NoErr(err)

	return ctx, cfg, be, dbx, dbstore, authTestKeys{
		user:           userKey,
		unknown:        unknownKey,
		bootstrapAdmin: adminKey,
	}
}

func mustGenerateKey(tb testing.TB, path string) (*keygen.KeyPair, gossh.PublicKey) {
	tb.Helper()

	is := is.New(tb)
	pair, err := keygen.New(path, keygen.WithKeyType(keygen.Ed25519), keygen.WithWrite())
	is.NoErr(err)

	pk, _, _, _, err := gossh.ParseAuthorizedKey([]byte(pair.AuthorizedKey()))
	is.NoErr(err)
	return pair, pk
}

// mockSSHContext implements ssh.Context for testing
type mockSSHContext struct {
	context.Context
	values      map[any]any
	permissions *ssh.Permissions
}

func newMockSSHContext(ctx context.Context) *mockSSHContext {
	return &mockSSHContext{
		Context:     ctx,
		values:      make(map[any]any),
		permissions: &ssh.Permissions{Permissions: &gossh.Permissions{Extensions: make(map[string]string)}},
	}
}

func (m *mockSSHContext) SetValue(key, value any) {
	m.values[key] = value
}

func (m *mockSSHContext) Value(key any) any {
	if v, ok := m.values[key]; ok {
		return v
	}
	return m.Context.Value(key)
}

func (m *mockSSHContext) Permissions() *ssh.Permissions {
	return m.permissions
}

func (m *mockSSHContext) User() string          { return "" }
func (m *mockSSHContext) RemoteAddr() net.Addr  { return &net.TCPAddr{} }
func (m *mockSSHContext) LocalAddr() net.Addr   { return &net.TCPAddr{} }
func (m *mockSSHContext) ServerVersion() string { return "" }
func (m *mockSSHContext) ClientVersion() string { return "" }
func (m *mockSSHContext) SessionID() string     { return "" }
func (m *mockSSHContext) Lock()                 {}
func (m *mockSSHContext) Unlock()               {}
