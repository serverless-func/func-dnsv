package main

import (
	"bytes"
	"fmt"
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
	args := []string{
		"-V",
	}
	cmd := exec.Command("markmap-cli", args...)

	// 捕获命令输出（可选）
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("命令执行失败: %v\n输出: %s", err, output)
	}
	return fmt.Sprintf("%s", output)
}
