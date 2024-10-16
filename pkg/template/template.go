package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type TemplateRaw struct {
	source string
	Tmpl   *Template
}

type Template struct {
	PM25 string
	PM10 string
	URL  string
}

func LoadTemplate(path string) (*TemplateRaw, error) {

	tmpl := &Template{}
	raw := &TemplateRaw{
		Tmpl: tmpl,
	}

	templateContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading template file: %s", err.Error())
	}

	raw.source = string(templateContent)
	return raw, nil
}

func (t *TemplateRaw) Parse() (string, error) {
	tmpl := template.New("weather template")
	tmpl, err := tmpl.Parse(t.source)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t.Tmpl)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

