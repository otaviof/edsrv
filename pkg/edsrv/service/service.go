package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/otaviof/edsrv/pkg/edsrv/editor"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Service represents the backend edit-server API, it handles the requests against the
// supported endpoints effectively exposing the application features.
type Service struct {
	logger *slog.Logger     // shared logger instance
	ed     editor.Interface // editor instance
}

const (
	// StatusPath status path.
	StatusPath = "/status"
	// StatusPath root path.
	RootPath = "/"

	// textPlain text-plain header.
	textPlain = "text/plain"
)

// status handles the requests for the "/status" endpoint, the response is based on local
// service configuration attributes.
func (s *Service) status(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType(textPlain)
	ctx.SetBodyString(fmt.Sprintf(
		"editor='%s', tmpDir='%s'", s.ed.GetCommand(), s.ed.GetTmpDir(),
	))
	ctx.SetStatusCode(http.StatusOK)

	s.logger.Debug("edit-server is running!", "endpoint", StatusPath)
}

// edit handles the requests for the "/" endpoint, the response is based on the Editor
// outcomes. The request body is informed for new file content, while the response body uses
// the final edited file payload.
func (s *Service) edit(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()
	logger := s.logger.With("endpoint", RootPath, "length", len(body))

	f, err := s.ed.Edit(body)
	if err != nil {
		logger.Error(err.Error())
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	logger = f.LoggerWith(logger)
	payload, err := f.Read()
	if err != nil {
		logger.Error(err.Error())
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("written", len(payload))
	logger.Debug("reading edited file")

	defer func() {
		if err := f.Remove(); err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Debug("temporary file removed")
	}()

	ctx.SetBody(payload)
	ctx.SetStatusCode(http.StatusOK)
	logger.Debug("all done!")
}

// RequestHandler instantiates the request router for the application endpoints.
func (s *Service) RequestHandler() fasthttp.RequestHandler {
	r := router.New()
	r.GET(StatusPath, s.status)
	r.POST(RootPath, s.edit)
	return r.Handler
}

// NewService returns a new service using a shared logger and editor instances.
func NewService(logger *slog.Logger, ed editor.Interface) *Service {
	return &Service{
		logger: logger,
		ed:     ed,
	}
}
