package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

// ErrNonSuccessfulStatusCode any other status-code than 200.
var ErrNonSuccessfulStatusCode = errors.New("non-successful status-code")

// StatusRequest executes a GET request on the edit-server status path, using the
// informed client address ("Addr" attribute) to create the request URI. Returns
// the response body and error when applicable.
func StatusRequest(logger *slog.Logger, c *fasthttp.HostClient) ([]byte, error) {
	statusURI, err := url.JoinPath("http://", c.Addr, StatusPath)
	if err != nil {
		return nil, err
	}
	logger.Info("dialing in...", "status-uri", statusURI)

	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.SetRequestURI(statusURI)
	req.Header.SetMethod(fasthttp.MethodGet)

	err = c.DoTimeout(req, res, 5*time.Second)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%w: %d",
			ErrNonSuccessfulStatusCode, res.StatusCode())
	}
	return res.Body(), nil
}
