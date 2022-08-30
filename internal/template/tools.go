package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

func echoReverseFunc(c echo.Context) func(string, ...interface{}) string {
	return func(resource string, params ...interface{}) string {
		return c.Echo().Reverse(resource, params...)
	}
}

func yieldFunc(tmpl *template.Template) func(name string, data interface{}) (string, error) {
	return func(name string, data interface{}) (string, error) {
		if t := tmpl.Lookup(name); t == nil {
			return "", nil
		}
		buf := bytes.NewBuffer([]byte{})
		err := tmpl.ExecuteTemplate(buf, name, data)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func global(globals map[string]interface{}) func(key string) (interface{}, error) {
	return func(key string) (interface{}, error) {
		value, ok := globals[key]
		if !ok {
			return nil, fmt.Errorf("global value for %q not found", key)
		}
		return value, nil
	}
}

func pathToID(path string) string {
	return strings.ReplaceAll(path, "/", "-")
}

func renderFuncs(c echo.Context, t *template.Template, globals map[string]interface{}) template.FuncMap {
	m := make(map[string]interface{}, 4)
	m["pathToID"] = pathToID
	m["lookupURI"] = echoReverseFunc(c)
	m["yield"] = yieldFunc(t)
	m["global"] = global(globals)

	return m
}
