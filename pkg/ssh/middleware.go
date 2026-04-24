package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"charm.land/log/v2"
	"charm.land/wish/v2"
	"github.com/charmbracelet/ssh"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/ssh/cmd"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/urutau-ltd/git-cone/pkg/store"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/time/rate"
)

// ErrPermissionDenied is returned when a user is not allowed connect.
var ErrPermissionDenied = fmt.Errorf("permission denied")

// AuthenticationMiddleware handles authentication.
func AuthenticationMiddleware(sh ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		// XXX: The authentication key is set in the context but gossh doesn't
		// validate the authentication. We need to verify that the _last_ key
		// that was approved is the one that's being used.

		ctx := s.Context()
		be := backend.FromContext(ctx)

		var pkFp string
		perms := s.Permissions().Permissions
		pk := s.PublicKey()
		if pk != nil {
			// There is no public key stored in the context, public-key auth
			// was never requested, skip
			if perms == nil {
				wish.Fatalln(s, ErrPermissionDenied)
				return
			}

			pkFp = gossh.FingerprintSHA256(pk)
		}

		// Check if the key is the same as the one we have in context
		fp := perms.Extensions["pubkey-fp"]
		if fp != "" && fp != pkFp {
			be.TrackAuthFailure(remoteIP(s))
			wish.Fatalln(s, ErrPermissionDenied)
			return
		}

		ac := be.AllowKeyless(ctx)
		publicKeyCounter.WithLabelValues(strconv.FormatBool(ac || pk != nil)).Inc()
		if !ac && pk == nil {
			be.TrackAuthFailure(remoteIP(s))
			wish.Fatalln(s, ErrPermissionDenied)
			return
		}

		// Set the auth'd user, or anon, in the context
		var user proto.User
		if pk != nil {
			user, _ = be.UserByPublicKey(ctx, pk)
		}
		ctx.SetValue(proto.ContextKeyUser, user)

		sh(s)
	}
}

func remoteIP(s ssh.Session) string {
	host, _, err := net.SplitHostPort(s.RemoteAddr().String())
	if err == nil {
		return host
	}
	return s.RemoteAddr().String()
}

// ContextMiddleware adds the config, backend, and logger to the session context.
func ContextMiddleware(cfg *config.Config, dbx *db.DB, datastore store.Store, be *backend.Backend, logger *log.Logger) func(ssh.Handler) ssh.Handler {
	return func(sh ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			ctx := s.Context()
			ctx.SetValue(sshutils.ContextKeySession, s)
			ctx.SetValue(config.ContextKey, cfg)
			ctx.SetValue(db.ContextKey, dbx)
			ctx.SetValue(store.ContextKey, datastore)
			ctx.SetValue(backend.ContextKey, be)
			ctx.SetValue(log.ContextKey, logger.WithPrefix("ssh"))
			sh(s)
		}
	}
}

var cliCommandCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "soft_serve",
	Subsystem: "cli",
	Name:      "commands_total",
	Help:      "Total times each command was called",
}, []string{"command"})

// CommandMiddleware handles git commands and CLI commands.
// This middleware must be run after the ContextMiddleware.
func CommandMiddleware(sh ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		_, _, ptyReq := s.Pty()
		if ptyReq {
			sh(s)
			return
		}

		ctx := s.Context()
		cfg := config.FromContext(ctx)

		args := s.Command()
		cliCommandCounter.WithLabelValues(cmd.CommandName(args)).Inc()
		rootCmd := &cobra.Command{
			Short:        "cone is a self-hosted Git server for the command line.",
			SilenceUsage: true,
		}
		rootCmd.CompletionOptions.DisableDefaultCmd = true

		rootCmd.SetUsageTemplate(cmd.UsageTemplate)
		rootCmd.SetUsageFunc(cmd.UsageFunc)
		rootCmd.AddCommand(
			cmd.GitUploadPackCommand(),
			cmd.GitUploadArchiveCommand(),
			cmd.GitReceivePackCommand(),
			cmd.RepoCommand(),
			cmd.SettingsCommand(),
			cmd.UserCommand(),
			cmd.InfoCommand(),
			cmd.PubkeyCommand(),
			cmd.SetUsernameCommand(),
			cmd.JWTCommand(),
			cmd.TokenCommand(),
			cmd.AuditCommand(),
		)

		if cfg.LFS.Enabled {
			rootCmd.AddCommand(
				cmd.GitLFSAuthenticateCommand(),
			)

			if cfg.LFS.SSHEnabled {
				rootCmd.AddCommand(
					cmd.GitLFSTransfer(),
				)
			}
		}

		rootCmd.SetArgs(args)
		if len(args) == 0 {
			// otherwise it'll default to os.Args, which is not what we want.
			rootCmd.SetArgs([]string{"--help"})
		}
		rootCmd.SetIn(s)
		rootCmd.SetOut(s)
		rootCmd.SetErr(s.Stderr())
		rootCmd.SetContext(ctx)

		if err := rootCmd.ExecuteContext(ctx); err != nil {
			s.Exit(1) //nolint: errcheck
			return
		}
	}
}

// LoggingMiddleware logs the ssh connection and command.
func LoggingMiddleware(sh ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		ctx := s.Context()
		logger := log.FromContext(ctx).WithPrefix("ssh")
		ct := time.Now()
		hpk := sshutils.MarshalAuthorizedKey(s.PublicKey())
		ptyReq, _, isPty := s.Pty()
		addr := s.RemoteAddr().String()
		user := proto.UserFromContext(ctx)
		logArgs := []interface{}{
			"addr",
			addr,
			"cmd",
			s.Command(),
		}

		if user != nil {
			logArgs = append([]interface{}{
				"username",
				user.Username(),
			}, logArgs...)
		}

		if isPty {
			logArgs = []interface{}{
				"term", ptyReq.Term,
				"width", ptyReq.Window.Width,
				"height", ptyReq.Window.Height,
			}
		}

		if config.IsVerbose() {
			logArgs = append(logArgs,
				"key", hpk,
				"envs", s.Environ(),
			)
		}

		msg := fmt.Sprintf("user %q", s.User())
		logger.Debug(msg+" connected", logArgs...)
		sh(s)
		logger.Debug(msg+" disconnected", append(logArgs, "duration", time.Since(ct))...)
	}
}

// InputRateLimitMiddleware limits SSH stdin events to prevent DoS via rapid
// key input (e.g. spamming Tab in the BubbleTea TUI causing CPU spikes).
func InputRateLimitMiddleware(sh ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		lim := rate.NewLimiter(rate.Every(50*time.Millisecond), 100)
		sh(limitedSession{Session: s, r: &rateLimitedReader{
			r:   s,
			lim: lim,
			ctx: s.Context(),
		}})
	}
}

type limitedSession struct {
	ssh.Session
	r io.Reader
}

func (ls limitedSession) Read(p []byte) (int, error) { return ls.r.Read(p) }

type rateLimitedReader struct {
	r   io.Reader
	lim *rate.Limiter
	ctx context.Context
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	if err := r.lim.Wait(r.ctx); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}
