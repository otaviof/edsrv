package service

import (
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/otaviof/edsrv/pkg/edsrv/editor"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/onsi/gomega"
)

//
// Test Snippet
// 	https://github.com/valyala/fasthttp/blob/dfce853067ec574be5311207da43d485711e0ca7/client_timing_test.go#L347-L389
//

const baseURI = "http://127.0.0.1/"

func TestService(t *testing.T) {
	g := gomega.NewWithT(t)

	payload := []byte("edited payload")
	ed := editor.NewFakeEditor(payload)
	s := NewService(slog.New(slog.Default().Handler()), ed)

	ln := fasthttputil.NewInmemoryListener()
	ch := make(chan struct{})
	go func() {
		if err := fasthttp.Serve(ln, s.RequestHandler()); err != nil {
			t.Errorf(err.Error())
		}
		close(ch)
	}()

	c := fasthttp.Client{Dial: func(_ string) (net.Conn, error) {
		return ln.Dial()
	}}

	t.Run(StatusPath, func(t *testing.T) {
		statusURL, err := url.JoinPath(baseURI, StatusPath)
		g.Expect(err).To(gomega.Succeed())
		t.Logf("status URL %q", statusURL)

		statusCode, body, err := c.Get(nil, statusURL)
		g.Expect(err).To(gomega.Succeed())
		g.Expect(statusCode).To(gomega.Equal(http.StatusOK))
		g.Expect(string(body)).To(gomega.ContainSubstring("fake-editor"))

		t.Logf("status request body %q", body)
	})

	t.Run(RootPath, func(t *testing.T) {
		t.Logf("root URL %q", baseURI)

		req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

		req.SetBody([]byte("payload"))
		req.SetRequestURI(baseURI)
		req.Header.SetMethod(fasthttp.MethodPost)

		err := c.DoTimeout(req, res, 5*time.Second)
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		g.Expect(err).To(gomega.Succeed())
		g.Expect(res.StatusCode()).To(gomega.Equal(http.StatusOK))
		g.Expect(res.Body()).To(gomega.Equal(payload))
	})

	ln.Close()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("server still running")
	}
}
