package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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

func Render(w func(text string)) (string, error) {
	w("DomainGroupWithRecords ...")
	gs := DomainGroupWithRecords(domain, w)
	w("DomainGroupWithRecords ... done")
	w("Template replace ...")
	var md bytes.Buffer
	err := tpl.Execute(&md, RenderRoot{
		Domain: domain,
		Groups: gs,
	})
	if err != nil {
		return "", err
	}
	w("Template replace ...")
	return md.String(), nil
}

func Convert(w func(text string)) string {
	mdContent, err := Render(w)
	if err != nil {
		return err.Error()
	}
	w("write result to md ...")
	err = os.WriteFile("/tmp/index.md", []byte(mdContent), 0644)
	if err != nil {
		return err.Error()
	}
	w("write result to md ... done")
	w("convert ...")
	cmd := exec.Command(
		"markmap",
		"--no-open",
		"--no-toolbar",
		"--output",
		"/tmp/index.html",
		"/tmp/index.md",
	)
	w("convert ... done")
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
	return strings.ReplaceAll(string(htmlContent[:]), "<title>Markmap</title>", "<title>DNS View</title>")
}
