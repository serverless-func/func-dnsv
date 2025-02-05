package main

import (
	"bytes"
	"os"
	"os/exec"
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

func Convert() string {
	mdContent, err := Render()
	if err != nil {
		return err.Error()
	}
	err = os.WriteFile("/tmp/index.md", []byte(mdContent), 0644)
	if err != nil {
		return err.Error()
	}
	cmd := exec.Command(
		"markmap",
		"--no-open",
		"--no-toolbar",
		"--output",
		"/tmp/index.html",
		"/tmp/index.md",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stderr.String()
	}
	htmlContent, err := os.ReadFile("/tmp/index.html")
	if err != nil {
		return err.Error()
	}
	return string(htmlContent[:])
}
