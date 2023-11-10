package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// EditBodyRequest sends a POST request to edit the informed payload through the
// edit-server, it returns the response body and error when applicable.
func EditBodyRequest(c *fasthttp.HostClient, body []byte) ([]byte, error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	req.SetBody(body)
	req.SetRequestURI(fmt.Sprintf("http://%s", c.Addr))
	req.Header.SetMethod(fasthttp.MethodPost)

	err := c.DoTimeout(req, res, 5*time.Second)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("non-successful status-code: %d", res.StatusCode())
	}
	return res.Body(), nil
}
