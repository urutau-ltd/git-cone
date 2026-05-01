package ssh

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"charm.land/log/v2"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/testsession"
	"github.com/charmbracelet/ssh"
	"github.com/matryer/is"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/db/migrate"
	"github.com/urutau-ltd/git-cone/pkg/store"
	"github.com/urutau-ltd/git-cone/pkg/store/database"
	gossh "golang.org/x/crypto/ssh"
	_ "modernc.org/sqlite" // sqlite driver
)

func TestSession(t *testing.T) {
	is := is.New(t)
	t.Run("authorized repo access", func(t *testing.T) {
		requireLocalListener(t)
		t.Log("setting up")
		s, close := setup(t)
		s.Stderr = os.Stderr
		t.Log("requesting pty")
		err := s.RequestPty("xterm", 80, 40, nil)
		is.NoErr(err)
		go func() {
			time.Sleep(1 * time.Second)
			// s.Signal(gossh.SIGTERM)
			s.Close() //nolint: errcheck
		}()
		t.Log("waiting for session to exit")
		_, err = s.Output("test")
		var ee *gossh.ExitMissingError
		is.True(errors.As(err, &ee))
		t.Log("session exited")
		is.NoErr(close())
	})
}

func requireLocalListener(tb testing.TB) {
	tb.Helper()

	lc := &net.ListenConfig{}
	l, err := lc.Listen(tb.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		tb.Skipf("local tcp listeners unavailable in this environment: %v", err)
		return
	}
	_ = l.Close()
}

func setup(tb testing.TB) (*gossh.Session, func() error) {
	tb.Helper()
	is := is.New(tb)
	dp := tb.TempDir()
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.DataPath = dp
	cfg.DB.DataSource = filepath.Join(dp, "test.db") + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	tb.Cleanup(func() {
		is.NoErr(os.RemoveAll(dp))
	})
	ctx = config.WithContext(ctx, cfg)
	dbx, err := db.Open(ctx, cfg.DB.Driver, cfg.DB.DataSource)
	if err != nil {
		tb.Fatal(err)
	}
	if err := migrate.Migrate(ctx, dbx); err != nil {
		tb.Fatal(err)
	}
	dbstore := database.New(ctx, dbx)
	ctx = store.WithContext(ctx, dbstore)
	be := backend.New(ctx, cfg, dbx, dbstore)
	return testsession.New(tb, &ssh.Server{
		Handler: ContextMiddleware(cfg, dbx, dbstore, be, log.Default())(bm.MiddlewareWithProgramHandler(SessionHandler)(func(s ssh.Session) {
			_, _, active := s.Pty()
			if !active {
				os.Exit(1)
			}
			s.Exit(0)
		})),
	}, nil), dbx.Close
}
