package main

import (
	"github.com/urutau-ltd/git-cone/internal/cli"
)

var (
	Version    = ""
	CommitSHA  = ""
	CommitDate = ""
)

func main() {
	app := cli.App{
		Use:        "cone",
		Short:      "A self-hosted Git server with an SSH-native TUI",
		Long:       "cone is a self-hosted Git server with an SSH-native TUI and a small operational footprint.",
		Copyright:  "(C) 2021-2026 Urutau Ltd.\nReleased under MIT license.",
		Container:  "cone",
		Version:    Version,
		CommitSHA:  CommitSHA,
		CommitDate: CommitDate,
	}
	app.Execute()
}
