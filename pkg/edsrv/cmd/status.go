package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/otaviof/edsrv/pkg/edsrv/config"
	"github.com/otaviof/edsrv/pkg/edsrv/service"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

// Status represents the "status" subcommand, which is meant to asses the application API
// backend current status.
type Status struct {
	logger *slog.Logger   // shared logger instance
	cmd    *cobra.Command // cobra instance
	cfg    *config.Config // flags for configuration
}

var statusDesc = fmt.Sprintf(`# %s status

Makes a noop request to assert edit-server status.

`, AppName)

// preRunE validates the application address.
func (s *Status) preRunE(_ *cobra.Command, _ []string) error {
	return s.cfg.ValidateAddrFlag()
}

// runE executes the request against the application status endpoint.
func (s *Status) runE(_ *cobra.Command, _ []string) error {
	statusURL, err := url.JoinPath("http://", s.cfg.Addr, service.StatusPath)
	if err != nil {
		return err
	}

	logger := s.cfg.LoggerWith(s.logger, config.AddrFlag).With("url", statusURL)
	logger.Info("dialing in...")

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(statusURL)
	req.Header.SetMethod(fasthttp.MethodGet)

	res := fasthttp.AcquireResponse()
	c := fasthttp.HostClient{Addr: s.cfg.Addr}

	err = c.DoTimeout(req, res, 15*time.Second)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	if err != nil {
		return fmt.Errorf("%w: connection error", err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("non-successful status code %d returned", res.StatusCode())
	}

	logger.Info(string(res.Body()))
	logger.Info("edit-server is healthy!")
	return nil
}

// NewStatus instantiate the "status" subcommand and its flags.
func NewStatus(logger *slog.Logger, cfg *config.Config) *Status {
	s := &Status{
		logger: logger,
		cmd: &cobra.Command{
			Use:          "status",
			Short:        "Checks edit-server status",
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
