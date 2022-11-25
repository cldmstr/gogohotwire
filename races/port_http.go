package races

import (
	"embed"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/ziflex/lecho/v3"

	"github.com/cldmstr/gogohotwire/internal/template"
)

const DomainPrefix = "races"

//go:embed views
var views embed.FS

type echoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type RaceUpdate struct {
	One   int
	Two   int
	Three int
}

type httpHandler struct {
	service  RacesService
	renderer *template.Renderer
	logger   *lecho.Logger
}

func RegisterRoutes(r echoRouter, t *template.Renderer, s RacesService, logger *lecho.Logger) error {
	h := httpHandler{
		service:  s,
		renderer: t,
		logger:   logger,
	}

	err := t.AddFS(DomainPrefix, views, false)
	if err != nil {
		return errors.Wrapf(err, "failed to register %s views", DomainPrefix)
	}

	r.GET("", h.list)
	r.GET("/:id/update", h.raceUpdate)
	r.GET("/:id", h.details)
	r.PUT("/:id/start", h.start)
	r.POST("", h.create)

	return nil
}

func (h *httpHandler) list(c echo.Context) error {
	races, err := h.service.Races()
	if err != nil {
		return errors.Wrap(err, "failed to load races")
	}

	data := map[string]interface{}{
		"Races": races,
	}

	return c.Render(http.StatusOK, "races/list.html", data)
}

func (h *httpHandler) details(c echo.Context) error {
	key := c.Param("id")
	id, err := uuid.Parse(key)
	if err != nil {
		return errors.Wrapf(err, "id not valid %q", key)
	}
	race, err := h.service.Race(id)
	if err != nil {
		return errors.Wrapf(err, "failed to load race %q", id.String())
	}

	templateName := "race_ready.partial.html"
	switch race.state {
	case Finished:
		templateName = "race_finished.partial.html"
	case Running:
		templateName = "race_running.partial.html"
	}

	data := map[string]interface{}{
		"Race": race,
	}

	yieldData := map[string]interface{}{
		"YieldTemplate": templateName,
		"YieldData":     data,
	}

	return c.Render(http.StatusOK, "races/race_details.stream.html", yieldData)
}

func (h *httpHandler) create(c echo.Context) error {
	raceName := c.FormValue("race-name")
	if raceName == "" {
		raceName = "404 Loop"
	}
	race, err := h.service.Create(raceName)
	if err != nil {
		return errors.Wrapf(err, "failed to create race with name %q", raceName)
	}

	data := map[string]interface{}{
		"Race": race,
	}

	return c.Render(http.StatusPermanentRedirect, "races/add_race.stream.html", data)
}

func (h *httpHandler) raceUpdate(c echo.Context) error {
	key := c.Param("id")
	id, err := uuid.Parse(key)
	if err != nil {
		return errors.Wrapf(err, "id not valid %q", key)
	}
	race, err := h.service.Race(id)

	if err != nil {
		return errors.Wrapf(err, "failed to load race %q", id.String())
	}
	ctx := c.Request().Context()
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().WriteHeader(http.StatusOK)

	if race.state != Running {
		return nil
	}

	updateChan := make(chan RaceUpdate, 1)
	defer close(updateChan)
	race.Watch(updateChan)

	select {
	case <-ctx.Done():
		return nil
	default:
		for r := range updateChan {
			err := h.renderer.RenderSSE(c.Response(), "races/race_update.html", map[string]interface{}{
				"Update": r,
			})
			if err != nil {
				return errors.Wrap(err, "failed to write update")
			}
		}
	}

	return nil
}

func (h *httpHandler) start(c echo.Context) error {
	key := c.Param("id")
	id, err := uuid.Parse(key)
	if err != nil {
		return errors.Wrapf(err, "failed to parse id %q", key)
	}
	race, err := h.service.Start(id)

	data := map[string]interface{}{
		"Race": race,
	}

	yieldData := map[string]interface{}{
		"YieldTemplate": "race_running.partial.html",
		"YieldData":     data,
	}

	return c.Render(http.StatusOK, "races/race_details.stream.html", yieldData)
}
