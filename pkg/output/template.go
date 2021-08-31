package output

import (
	"bytes"
	"text/template"
)

func FromTemplate(tmpl string, ob interface{}) string {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return ""
	}
	var out bytes.Buffer
	err = t.Execute(&out, ob)
	if err != nil {
		return ""
	}
	return out.String()
}
