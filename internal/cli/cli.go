package cli

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"charm.land/log/v2"
	"github.com/charmbracelet/colorprofile"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/cmd/soft/admin"
	"github.com/urutau-ltd/git-cone/cmd/soft/browse"
	"github.com/urutau-ltd/git-cone/cmd/soft/hook"
	"github.com/urutau-ltd/git-cone/cmd/soft/serve"
	"github.com/urutau-ltd/git-cone/pkg/config"
	logr "github.com/urutau-ltd/git-cone/pkg/log"
	"github.com/urutau-ltd/git-cone/pkg/ui/common"
	"github.com/urutau-ltd/git-cone/pkg/version"
	"go.uber.org/automaxprocs/maxprocs"
)

type App struct {
	Use       string
	Short     string
	Long      string
	Copyright string
	Container string

	Version    string
	CommitSHA  string
	CommitDate string
}

func (a *App) Execute() {
	rootCmd := &cobra.Command{
		Use:          a.Use,
		Short:        a.Short,
		Long:         a.Long,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return browse.Command.RunE(cmd, args)
		},
	}

	manCmd := &cobra.Command{
		Use:    "man",
		Short:  "Generate man pages",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			manPage, err := mcobra.NewManPage(1, rootCmd)
			if err != nil {
				return err
			}

			manPage = manPage.WithSection("Copyright", a.Copyright)
			fmt.Println(manPage.Build(roff.NewDocument()))
			return nil
		},
	}

	noColor, _ := strconv.ParseBool(os.Getenv("GIT_CONE_NO_COLOR"))
	if !noColor {
		noColor, _ = strconv.ParseBool(os.Getenv("SOFT_SERVE_NO_COLOR"))
	}
	if noColor {
		common.DefaultColorProfile = colorprofile.NoTTY
	}

	rootCmd.AddCommand(
		manCmd,
		serve.Command,
		hook.Command,
		admin.Command,
		browse.Command,
	)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	if len(a.CommitSHA) >= 7 {
		vt := rootCmd.VersionTemplate()
		rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + a.CommitSHA[0:7] + ")\n")
	}
	if a.Version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			a.Version = info.Main.Version
		} else {
			a.Version = "unknown (built from source)"
		}
	}
	rootCmd.Version = a.Version

	version.Version = a.Version
	version.CommitSHA = a.CommitSHA
	version.CommitDate = a.CommitDate

	ctx := context.Background()
	cfg := config.DefaultConfig()
	if cfg.Exist() {
		if err := cfg.Parse(); err != nil {
			log.Fatal(err)
		}
	}

	if err := cfg.ParseEnv(); err != nil {
		log.Fatal(err)
	}

	ctx = config.WithContext(ctx, cfg)
	logger, f, err := logr.NewLogger(cfg)
	if err != nil {
		log.Errorf("failed to create logger: %v", err)
	}

	ctx = log.WithContext(ctx, logger)
	if f != nil {
		defer f.Close() //nolint: errcheck
	}

	log.SetDefault(logger)

	var opts []maxprocs.Option
	if config.IsVerbose() {
		opts = append(opts, maxprocs.Logger(log.Debugf))
	}

	if _, err := maxprocs.Set(opts...); err != nil {
		log.Warn("couldn't set automaxprocs", "error", err)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
