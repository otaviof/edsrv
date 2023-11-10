package service

import (
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/otaviof/edsrv/pkg/edsrv/editor"
	"github.com/otaviof/edsrv/test/helper"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	g := NewWithT(t)

	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	))

	payload := []byte("edited payload")
	ed := editor.NewFakeEditor(payload)
	srv := NewService(logger, ed)

	ln := fasthttputil.NewInmemoryListener()
	ch := make(chan struct{})
	go func() {
		if err := fasthttp.Serve(ln, srv.RequestHandler()); err != nil {
			t.Errorf(err.Error())
		}
		close(ch)
	}()

	c := &fasthttp.HostClient{
		Addr: "127.0.0.1:1982",
		Dial: func(_ string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	t.Run(StatusPath, func(_ *testing.T) {
		resBody, err := StatusRequest(logger, c)
		g.Expect(err).To(Succeed())
		g.Expect(string(resBody)).To(ContainSubstring(ed.GetCommand()))
		g.Expect(string(resBody)).To(ContainSubstring(ed.GetTmpDir()))

		t.Logf("status request body %q", resBody)
	})

	t.Run(RootPath, func(_ *testing.T) {
		resBody, err := helper.EditBodyRequest(c, []byte("initial input..."))
		g.Expect(err).To(Succeed())
		g.Expect(resBody).To(Equal(payload))
	})

	ln.Close()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("server still running")
	}
}
