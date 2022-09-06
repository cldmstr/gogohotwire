package template

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

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

func renderFuncs(t *template.Template) template.FuncMap {
	m := sprig.TxtFuncMap()
	m["yield"] = yieldFunc(t)

	return m
}
