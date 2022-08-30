package races

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/cldmstr/gogohotwire/internal/template"
)

const DomainPrefix = "races"

//go:embed views
var views embed.FS

type echoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type httpHandler struct {
}

func RegisterRoutes(r echoRouter, t *template.Renderer) error {
	h := httpHandler{}

	err := t.AddFS(DomainPrefix, views, false)
	if err != nil {
		return errors.Wrapf(err, "failed to register %s views", DomainPrefix)
	}

	r.GET("", h.list)
	r.POST("", h.new)

	return nil
}

func (h *httpHandler) list(c echo.Context) error {

	data := map[string]interface{}{
		"Races": []string{"Prairie Circuit", "South Rodent Ring", "Gophtona 500"},
	}

	return c.Render(http.StatusOK, "races/list.html", data)
}

func (h *httpHandler) new(c echo.Context) error {

	raceName := c.Request().Form.Get("race-name")
	if raceName == "" {
		raceName = "404 Loop"
	}

	data := map[string]interface{}{
		"Race": raceName,
	}

	return c.Render(http.StatusOK, "races/add_race.stream.html", data)
}
