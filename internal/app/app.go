package app

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/cldmstr/gogohotwire/internal/races"

	"github.com/cldmstr/gogohotwire/internal/drivers"

	"github.com/cldmstr/gogohotwire/assets"
	"github.com/cldmstr/gogohotwire/internal/template"
)

//go:embed views
var views embed.FS

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	e := echo.New()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerologger := zerolog.New(zerolog.NewConsoleWriter())
	zerologger = zerologger.Level(zerolog.DebugLevel)

	logger := lecho.From(zerologger,
		lecho.WithTimestamp(),
		lecho.WithPrefix("app"),
	)
	e.Logger = logger
	e.Use(
		lecho.Middleware(lecho.Config{Logger: logger}))

	renderer, err := template.New(false)
	if err != nil {
		return errors.Wrap(err, "failed to create template renderer")
	}
	err = renderer.AddFS("app", views, true)
	if err != nil {
		return errors.Wrap(err, "failed to setup app views")
	}

	e.Renderer = renderer
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/static/*", func(c echo.Context) error {
		// if a request in /static is not handled by the static middleware,
		// the requested resource does not exist.
		return errors.Errorf("static resource not found %q", c.Request().URL.Path)
	},
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "static",
			Browse:     false,
			Filesystem: http.FS(assets.Assets),
		}))
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/static/images/favicon.ico")
	})
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/races")
	})

	racesService := races.New()

	racesGroup := e.Group(races.DomainPrefix)
	err = races.RegisterRoutes(racesGroup, renderer, racesService, logger)
	if err != nil {
		return errors.Wrap(err, "failed to register races domain")
	}

	driversService := drivers.NewDriversService(drivers.NewDriversStore())
	driversGroup := e.Group(drivers.DomainPrefix)
	err = drivers.RegisterRoutes(driversGroup, renderer, driversService)
	if err != nil {
		return errors.Wrap(err, "failed to register drivers domain")
	}

	return e.Start("localhost:8088")
}
