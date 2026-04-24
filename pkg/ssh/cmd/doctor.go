package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/sshpolicy"
	"github.com/urutau-ltd/git-cone/pkg/version"
)

func DoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "doctor",
		Short:             "Show effective server security and runtime settings",
		Args:              cobra.NoArgs,
		PersistentPreRunE: checkIfAdmin,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cfg := config.FromContext(ctx)

			ver := version.Version
			if ver == "" {
				ver = "unknown"
			}

			cmd.Printf("server:           git-cone %s\n", ver)
			cmd.Printf("strict:           %s\n", yesNo(cfg.Security.Strict))
			cmd.Printf("ssh-listen:       %s\n", cfg.SSH.ListenAddr)
			cmd.Printf("ssh-public-url:   %s\n", cfg.SSH.PublicURL)
			cmd.Printf("http-listen:      %s\n", cfg.HTTP.ListenAddr)
			cmd.Printf("http-public-url:  %s\n", cfg.HTTP.PublicURL)
			cmd.Printf("stats-listen:     %s\n", cfg.Stats.ListenAddr)
			cmd.Printf("git-listen:       %s\n", stringOrDisabled(cfg.Git.ListenAddr))
			cmd.Printf("lfs-enabled:      %s\n", yesNo(cfg.LFS.Enabled))
			cmd.Printf("lfs-ssh-enabled:  %s\n", yesNo(cfg.LFS.SSHEnabled))
			cmd.Printf("hook-timeout:     %ds\n", cfg.Hooks.Timeout)
			cmd.Printf("ssh-idle-timeout: %ds\n", cfg.SSH.IdleTimeout)
			cmd.Printf("ssh-max-timeout:  %ds\n", cfg.SSH.MaxTimeout)
			cmd.Printf("host-key-path:    %s\n", cfg.SSH.KeyPath)
			cmd.Printf("host-key-exists:  %s\n", yesNo(pathExists(cfg.SSH.KeyPath)))
			cmd.Printf("client-key-path:  %s\n", cfg.SSH.ClientKeyPath)
			cmd.Printf("client-key-exists: %s\n", yesNo(pathExists(cfg.SSH.ClientKeyPath)))
			cmd.Printf("kex:              %s\n", strings.Join(sshpolicy.HardenedKeyExchanges(), ", "))
			cmd.Printf("ciphers:          %s\n", strings.Join(sshpolicy.HardenedCiphers(), ", "))
			cmd.Printf("macs:             %s\n", strings.Join(sshpolicy.HardenedMACs(), ", "))
			return nil
		},
	}
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func stringOrDisabled(value string) string {
	if value == "" {
		return "disabled"
	}
	return value
}
