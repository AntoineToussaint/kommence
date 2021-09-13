package output

import (
	"bytes"
	"text/template"
)

func FromTemplate(log *Logger, tmpl string, ob interface{}) string {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		log.Errorf("can't load template: %v", err)
		return ""
	}
	var out bytes.Buffer
	err = t.Execute(&out, ob)
	if err != nil {
		log.Errorf("can't templatize: %v", err)
		return ""
	}
	return out.String()
}
