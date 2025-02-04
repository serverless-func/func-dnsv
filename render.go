package main

import (
	"bytes"
	"text/template"
)

type RenderRoot struct {
	Domain string
	Groups []RecordGroup
}

var domain = "dongfg.com"

var tpl = template.Must(template.New("markdown").Parse(`---
title: {{ .Domain }}
markmap:
  colorFreezeLevel: 2
---

{{- range .Groups}}
## {{.GroupName}}

{{range .Records -}}
- {{ .Name }} {{ .Remark }}
{{end}}
{{- end}}`))

func Render() (string, error) {
	gs := DomainGroupWithRecords(domain)
	var md bytes.Buffer
	err := tpl.Execute(&md, RenderRoot{
		Domain: domain,
		Groups: gs,
	})
	if err != nil {
		return "", err
	}
	return md.String(), nil
}
