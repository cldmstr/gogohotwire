package drivers

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/cldmstr/gogohotwire/internal/template"
)

const DomainPrefix = "drivers"

//go:embed views
var views embed.FS

type echoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type htmlHandler struct {
	service DriversService
}

func RegisterRoutes(r echoRouter, t *template.Renderer, service DriversService) error {
	h := htmlHandler{
		service: service,
	}

	err := t.AddFS(DomainPrefix, views, false)
	if err != nil {
		return errors.Wrapf(err, "failed to register %s views", DomainPrefix)
	}

	r.GET("", h.list)

	return nil
}

func (h *htmlHandler) list(c echo.Context) error {
	drivers, err := h.service.Drivers()
	if err != nil {
		return errors.Wrapf(err, "failed to load drivers")
	}

	data := map[string]interface{}{
		"Drivers": drivers,
	}

	return c.Render(http.StatusOK, "drivers/card_list.html", data)
}
