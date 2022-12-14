package template

import (
	"bytes"
	"embed"
	"encoding/base64"
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

func (r *Renderer) AddFS(namespace string, fsys fs.FS, isGlobal bool) error {
	sub, err := fs.Sub(fsys, "views")
	if err != nil {
		return errors.Wrapf(err, "failed to add filesystem %q", namespace)
	}
	r.templates[namespace] = sub
	if isGlobal {
		r.globalPaths = append(r.globalPaths, filepath.Join(namespace, "*.html"))
	}

	return nil
}

func (r *Renderer) RenderSSE(res *echo.Response, name string, data interface{}) error {
	tmpl, err := r.newTemplate(name)
	if err != nil {
		return errors.Wrap(err, "failed to initialize template")
	}

	res.Header().Set(echo.HeaderContentType, "text/event-stream")

	event := bytes.NewBufferString("data: ")
	encoder := base64.NewEncoder(base64.StdEncoding, event)
	defer encoder.Close()

	err = tmpl.ExecuteTemplate(encoder, filepath.Base(name), data)
	if err != nil {
		return errors.Wrapf(err, "failed to template %q", name)
	}
	event.WriteString("\n\n")

	_, err = res.Write(event.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to write event")
	}
	res.Flush()

	return nil
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := r.newTemplate(name)
	if err != nil {
		return errors.Wrap(err, "failed to initialize template")
	}

	applicationData := map[string]interface{}{
		"YieldTemplate": filepath.Base(name),
		"YieldData":     data,
	}

	switch {
	case strings.HasSuffix(name, ".stream.html"):
		c.Response().Header().Add("Content-Type", "text/vnd.turbo-stream.html")
		err = tmpl.ExecuteTemplate(w, filepath.Base(name), data)
	case c.Request().Header.Get("Turbo-Frame") != "":
		err = tmpl.ExecuteTemplate(w, filepath.Base(name), data)
	default:
		err = tmpl.ExecuteTemplate(w, "application.html", applicationData)
	}
	if err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}

func (r *Renderer) newTemplate(name string) (*template.Template, error) {
	tmpl := template.New("")
	tmpl.Funcs(renderFuncs(tmpl))
	namespace := filepath.Join(filepath.Dir(name), "*.html")
	tmpl, err := tmpl.ParseFS(r.templates, append(r.globalPaths, namespace)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load templates %q", namespace)
	}

	return tmpl, nil
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
