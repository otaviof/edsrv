package integration_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/otaviof/edsrv/pkg/edsrv/cmd"
	"github.com/otaviof/edsrv/pkg/edsrv/config"
	"github.com/otaviof/edsrv/test/helper"

	"github.com/valyala/fasthttp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("edsrv integration tests", Ordered, func() {
	// contains the shared logger instance
	var logger *slog.Logger

	// contains the value for the flags used on start and status subcommands
	var addrFlagValue string
	var editorFlagValue string
	var tmpDirFlagValue string

	// cancel function to stop "edsrv start" background process
	var cancelFn context.CancelFunc

	// instantitate the shared variables that will be used on the test-cases below
	BeforeAll(func() {
		slogOpts := &slog.HandlerOptions{Level: slog.LevelDebug}
		logger = slog.New(slog.NewTextHandler(GinkgoWriter, slogOpts))

		addrFlagValue = getAddrFlagFromEnvOrDefault()

		tmpFile, err := createStaticFileOnTmpDir()
		Expect(err).To(Succeed())

		tmpDirFlagValue = filepath.Dir(tmpFile)
		editorFlagValue = fmt.Sprintf("cp -f -v %s", tmpFile)
	})

	// stops the edit-server running in the background and removes the temporary
	// directory used during testing
	AfterAll(func() {
		cancelFn()
		_ = os.RemoveAll(tmpDirFlagValue)
	})

	// describes the regular usage workflow, the user will first start the
	// edit-server in the background and then try to run "edsrv status", then edit
	// a file
	Context("edsrv regular workflow", func() {
		It("CMD: 'edsrv start' (running in the background)", func() {
			startCmd := cmd.NewStart(logger, config.NewConfig()).Cmd()

			err := startCmd.ParseFlags(stringMapToSlice(map[string]string{
				fmt.Sprintf("--%s", config.AddrFlag):   addrFlagValue,
				fmt.Sprintf("--%s", config.EditorFlag): editorFlagValue,
				fmt.Sprintf("--%s", config.TmpDirFlag): tmpDirFlagValue,
			}))
			Expect(err).To(Succeed())

			err = startCmd.PreRunE(startCmd, startCmd.Flags().Args())
			Expect(err).To(Succeed())

			// registering a context and extracting a cancel function for the
			// "edsrv start" running in the background
			var ctx context.Context
			ctx, cancelFn = context.WithCancel(context.Background())
			startCmd.SetContext(ctx)

			go func() {
				defer GinkgoRecover()
				err := startCmd.RunE(startCmd, startCmd.Flags().Args())
				Expect(err).To(Succeed())
			}()
		})

		It("CMD: 'edsrv status'", func() {
			statusCmd := cmd.NewStatus(logger, config.NewConfig()).Cmd()

			err := statusCmd.ParseFlags(stringMapToSlice(map[string]string{
				fmt.Sprintf("--%s", config.AddrFlag): addrFlagValue,
			}))
			Expect(err).To(Succeed())

			err = statusCmd.PreRunE(statusCmd, statusCmd.Flags().Args())
			Expect(err).To(Succeed())

			// retrying "edsrv status" a few times in order to wait for the
			// edit-server starting up, normally this works right on the first
			// attempt
			Eventually(func() error {
				return statusCmd.RunE(statusCmd, statusCmd.Flags().Args())
			}).WithTimeout(5 * time.Second).Should(Succeed())
		})

		// simulates what a browser extension does, an arbitrary payload is edited
		// through the edit-server API, the result must match what's expected
		It("/POST: edit request payload", func() {
			c := &fasthttp.HostClient{Addr: addrFlagValue}
			resBody, err := helper.EditBodyRequest(c, []byte("initial input"))
			Expect(err).To(Succeed())
			Expect(resBody).To(Equal(staticFilePayload))
		})
	})
})
