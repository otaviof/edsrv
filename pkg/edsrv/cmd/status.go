package cmd

import (
	"fmt"
	"log/slog"

	"github.com/otaviof/edsrv/pkg/edsrv/config"
	"github.com/otaviof/edsrv/pkg/edsrv/service"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

// Status represents the "status" subcommand, which is meant to asses the
// application API backend current status.
type Status struct {
	logger *slog.Logger   // shared logger instance
	cmd    *cobra.Command // cobra instance
	cfg    *config.Config // flags for configuration
}

var statusDesc = fmt.Sprintf(`# %s status

Makes a noop request to assert edit-server status.

`, AppName)

// Cmd exposes the cobra command instance.
func (s *Status) Cmd() *cobra.Command {
	return s.cmd
}

// preRunE validates the application address.
func (s *Status) preRunE(_ *cobra.Command, _ []string) error {
	return s.cfg.ValidateAddrFlag()
}

// runE executes the request against the application status endpoint.
func (s *Status) runE(_ *cobra.Command, _ []string) error {
	logger := s.cfg.LoggerWith(s.logger, config.AddrFlag)

	c := &fasthttp.HostClient{Addr: s.cfg.Addr}
	resBody, err := service.StatusRequest(logger, c)
	if err != nil {
		return err
	}

	logger.Info(string(resBody))
	logger.Info("edit-server is healthy!")
	return nil
}

// NewStatus instantiates the "status" subcommand and its flags.
func NewStatus(logger *slog.Logger, cfg *config.Config) *Status {
	s := &Status{
		logger: logger,
		cmd: &cobra.Command{
			Use:          "status",
			Short:        "Verifies the edit-server status",
			Long:         statusDesc,
			SilenceUsage: true,
		},
		cfg: cfg,
	}
	s.cmd.PreRunE = s.preRunE
	s.cmd.RunE = s.runE
	s.cfg.AddAddrFlag(s.cmd.PersistentFlags())
	return s
}
