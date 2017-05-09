package main

import (
	"text/template"
)

var tmpl = template.Must(template.New("code").Parse(`{{ range $k, $v := . }}
	{{- $k }}: {{ $v }}
{{ end }}`))
