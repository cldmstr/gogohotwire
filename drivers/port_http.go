package drivers

import (
	"embed"
	"net/http"
	"strconv"

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
	r.GET("/:id/icon", h.icon)
	r.GET("/:id/details", h.details)

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

	return c.Render(http.StatusOK, "drivers/list.html", data)
}

func (h *htmlHandler) icon(c echo.Context) error {
	key := c.Param("id")
	id, err := strconv.Atoi(key)
	if err != nil {
		return errors.Wrapf(err, "failed due to baddly formatted id %q", key)
	}

	driver, err := h.service.Driver(id)
	if err != nil {
		return errors.Errorf("failed to load driver %q", key)
	}

	return c.Render(http.StatusOK, "drivers/icon.html", driver)
}

func (h *htmlHandler) details(c echo.Context) error {
	key := c.Param("id")
	id, err := strconv.Atoi(key)
	if err != nil {
		return errors.Wrapf(err, "failed due to baddly formatted id %q", key)
	}

	driver, err := h.service.Driver(id)
	if err != nil {
		return errors.Errorf("failed to load driver %q", key)
	}

	return c.Render(http.StatusOK, "drivers/details.html", driver)
}
