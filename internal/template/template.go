package template

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

//go:embed views
var views embed.FS

type Renderer struct {
	templates       templateFS
	globalPaths     []string
	developmentMode bool
}

// New setups a new template renderer.
func New(developmentMode bool) (*Renderer, error) {
	t := &Renderer{
		templates:       make(map[string]fs.FS),
		developmentMode: developmentMode,
	}

	err := t.AddFS("tmpl", views, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load general templates")
	}

	return t, nil
}

func (t *Renderer) AddFS(namespace string, fsys fs.FS, isGlobal bool) error {
	sub, err := fs.Sub(fsys, "views")
	if err != nil {
		return errors.Wrapf(err, "failed to add filesystem %q", namespace)
	}
	t.templates[namespace] = sub
	if isGlobal {
		t.globalPaths = append(t.globalPaths, filepath.Join(namespace, "*.html"))
	}

	return nil
}

func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	globals := map[string]interface{}{
		"csrfToken":   c.Get("csrf"),
		"print":       c.QueryParam("print") != "",
		"development": t.developmentMode,
	}

	tmpl := template.New("")
	tmpl.Funcs(renderFuncs(c, tmpl, globals))
	namespace := filepath.Join(filepath.Dir(name), "*.html")
	tmpl, err := tmpl.ParseFS(t.templates, append(t.globalPaths, namespace)...)
	if err != nil {
		return errors.Wrapf(err, "failed to load templates %q", namespace)
	}

	applicationData := map[string]interface{}{
		"YieldTemplate": filepath.Base(name),
		"YieldData":     data,
		"Print":         c.QueryParam("print") != "",
	}

	switch {
	case strings.HasSuffix(name, ".stream.html"):
		c.Response().Header().Add("Content-Type", "text/vnd.turbo-stream.html")
		err = tmpl.ExecuteTemplate(w, filepath.Base(name), data)
	case c.Request().Header.Get("Turbo-Frame") != "":
		applicationData["FrameID"] = c.Request().Header.Get("Turbo-Frame")
		err = tmpl.ExecuteTemplate(w, "turbo-frame.html", applicationData)
	default:
		err = tmpl.ExecuteTemplate(w, "application.html", applicationData)
	}
	if err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}

type templateFS map[string]fs.FS

var (
	_ fs.FS         = &templateFS{}
	_ fs.GlobFS     = &templateFS{}
	_ fs.ReadFileFS = &templateFS{}
)

func (t templateFS) ReadFile(name string) ([]byte, error) {
	ns, internalName, err := t.internalName(name)
	if err != nil {
		return nil, err
	}
	fsys, err := t.getFS(ns)
	if err != nil {
		return nil, err
	}

	return fs.ReadFile(fsys, internalName)
}

func (t templateFS) Open(name string) (fs.File, error) {
	ns, internalName, err := t.internalName(name)
	if err != nil {
		return nil, err
	}
	fsys, err := t.getFS(ns)
	if err != nil {
		return nil, err
	}

	return fsys.Open(internalName)
}

func (t templateFS) Glob(pattern string) ([]string, error) {
	namespace, internalPattern, err := t.internalName(pattern)
	if err != nil {
		return nil, err
	}
	fsys, err := t.getFS(namespace)
	if err != nil {
		return nil, err
	}

	matches, err := fs.Glob(fsys, internalPattern)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(matches))
	for _, n := range matches {
		paths = append(paths, filepath.Join(namespace, n))
	}
	return paths, nil
}

func (t templateFS) internalName(name string) (namespace string, internalName string, err error) {
	p := strings.Split(name, string(filepath.Separator))
	if len(p) < 2 {
		return "", "", errors.New("path requires at least two levels")
	}
	namespace = p[0]
	internalName = filepath.Join(p[1:]...)

	return namespace, internalName, err
}

func (t templateFS) getFS(namespace string) (fs.FS, error) {
	if fsys, ok := t[namespace]; ok {
		return fsys, nil
	}

	return nil, fs.ErrNotExist
}
