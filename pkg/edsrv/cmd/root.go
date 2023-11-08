package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/otaviof/edsrv/pkg/edsrv/config"

	"github.com/spf13/cobra"
)

// Root represents the primary application command.
type Root struct {
	cmd *cobra.Command // cobra instance
	cfg *config.Config // shared configuration
}

// AppName application name.
const AppName = "edsrv"

var rootDesc = fmt.Sprintf(`# %s

Is a edit-server meant to work as a API backend for browser extensions that allow using a
external text editor on a regular webpage.

`, AppName)

// Cmd shares the cobra.Command instance decorated with subcommands.
func (r *Root) Cmd() *cobra.Command {
	logOpts := &slog.HandlerOptions{Level: r.cfg.LogLevel}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))

	r.cmd.AddCommand(NewStart(logger, r.cfg).cmd)
	r.cmd.AddCommand(NewStatus(logger, r.cfg).cmd)

	return r.cmd
}

// NewRoot instantiate the root command with global configuration instance and flags.
func NewRoot() *Root {
	r := &Root{
		cmd: &cobra.Command{
			Use:   AppName,
			Short: "backend for browser extensions to use a external text editors",
			Long:  rootDesc,
		},
		cfg: config.NewConfig(),
	}
	r.cfg.AddLogLevelFlag(r.cmd.PersistentFlags())
	return r
}
