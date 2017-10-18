package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"text/template"
)

var tmpl = template.New("code")

func init() {
	token := func(t Token) string {
		switch t := t.(type) {
		case Term:
			return fmt.Sprintf("newTerm(%q)", t)

		case NTerm:
			return fmt.Sprintf("newNTerm(%q)", t[1:len(t)-1])

		case Epsilon:
			return "newEpsilon()"

		case EOF:
			return "newEOF()"
		}

		panic(fmt.Errorf("Unexpected token type: %T", t))
	}

	tmpl.Funcs(template.FuncMap{
		"isEpsilon": isEpsilon,
		"token":     token,
		"rule": func(r Rule) string {
			var buf bytes.Buffer
			buf.WriteString("newRule(")

			var s string
			for _, t := range r {
				buf.WriteString(s)
				buf.WriteString(token(t))
				s = ", "
			}

			buf.WriteRune(')')
			return buf.String()
		},
	})

	tmpl = template.Must(tmpl.Parse(`// This file was auto-generated by pgen.
// Editing is not advised.

package pgen

var Table = map[Lookup]Rule{
	{{ range $nterm, $_ := . }}
		{{- range $term, $from := $.First $nterm -}}
			{{- if isEpsilon $term | not -}}
				{Term: {{ $term | token }}, NTerm: {{ $nterm | token -}} }: {{ $from | rule }},
			{{ end -}}
		{{ end -}}

		{{- if $.Nullable $nterm -}}
		{{- range $term, $from := $.Follow $nterm -}}
			{{- if isEpsilon $term | not -}}
				{Term: {{ $term | token }}, NTerm: {{ $nterm | token -}} }: newRule(newEpsilon()),
			{{ end -}}
		{{ end -}}
		{{ end -}}
	{{ end }}
}`))
}

type formatter struct {
	w   io.Writer
	buf bytes.Buffer
}

func (f *formatter) Write(data []byte) (int, error) {
	return f.buf.Write(data)
}

func (f formatter) Close() error {
	src, err := format.Source(f.buf.Bytes())
	if err != nil {
		return err
	}

	_, err = f.w.Write(src)
	return err
}
