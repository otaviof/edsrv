package cmd

import (
	"fmt"
	"log/slog"

	"github.com/otaviof/edsrv/pkg/edsrv/config"
	"github.com/otaviof/edsrv/pkg/edsrv/editor"
	"github.com/otaviof/edsrv/pkg/edsrv/service"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

// Start represents the "start" subcommand, which starts the application API
// backend server.
type Start struct {
	logger *slog.Logger   // shared logger instance
	cmd    *cobra.Command // cobra instance
	cfg    *config.Config // flags for configuration
}

var startDesc = fmt.Sprintf(`# %s start

Starts the edit-server API backend using the informed flags for configuration.

`, AppName)

// Cmd exposes the cobra command instance.
func (s *Start) Cmd() *cobra.Command {
	return s.cmd
}

// preRunE Validates the informed configuration.
func (s *Start) preRunE(_ *cobra.Command, _ []string) error {
	return s.cfg.ValidateStartFlags()
}

// runE runs the API backend using the configuration informed via flags.
func (s *Start) runE(_ *cobra.Command, _ []string) error {
	logger := s.cfg.LoggerWith(
		s.logger, config.AddrFlag, config.EditorFlag, config.TmpDirFlag,
	)

	ed := editor.NewEditor(logger, s.cfg.Editor, s.cfg.TmpDir)
	srv := service.NewService(logger, ed)

	logger.Debug("starting edit-server...")
	return fasthttp.ListenAndServe(s.cfg.Addr, srv.RequestHandler())
}

// NewStart instantiates "start" subcommand and its flags.
func NewStart(logger *slog.Logger, cfg *config.Config) *Start {
	s := &Start{
		logger: logger,
		cmd: &cobra.Command{
			Use:          "start",
			Short:        "Starts the edit-server API backend service",
			Long:         startDesc,
			SilenceUsage: true,
		},
		cfg: cfg,
	}
	s.cmd.PreRunE = s.preRunE
	s.cmd.RunE = s.runE
	s.cfg.AddStartFlags(s.cmd.PersistentFlags())
	return s
}
