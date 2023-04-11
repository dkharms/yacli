package yacli

import (
	"fmt"
	"text/template"
)

var (
	formatRed     = func(msg string) string { return fmt.Sprintf("\033[31m%s\033[0m", msg) }
	formatBlue    = func(msg string) string { return fmt.Sprintf("\033[34m%s\033[0m", msg) }
	formatDefault = func(msg string) string { return fmt.Sprintf("\033[0m%s\033[0m", msg) }
)

var (
	formatBold = func(msg string) string { return fmt.Sprintf("\033[1m%s\033[0m", msg) }
)

var helpTemplateRaw = `{{ if .Deprecated }}[{{ FormatRed "DEPRECATED" }}] {{ end }}{{ .Usage }}
{{ .Description }}

Flags:
{{- range .Flags }}
    {{ if .Deprecated }}[{{ FormatRed "DEPRECATED" }}] {{ end }}{{ printf "-%s" .Short | FormatBold }} | {{printf "--%s" .Name | FormatBold }} [{{ printf "%s" .Type | FormatBlue }}] - {{ .Description }} 
{{- end }}
{{- if gt (len .Arguments) 0 }} 

Arguments:
{{- range .Arguments }}
    {{ if not .Optional}}{{ FormatBold "*" | FormatRed }}{{end}} {{ FormatBold .Name }} [{{ printf "%s" .Type | FormatBlue }}] - {{ .Description }} 
{{- end }}
{{- end }}
{{- if gt (len .Subcommands) 0 }}

Subcommands:
{{- range .Subcommands }}
    {{ FormatBold .Name }} - {{ .Description }} {{ if .Deprecated }}[{{ FormatRed "DEPRECATED" }}]{{ end }}
{{- end }}
{{- end }}
`

var helpTemplate = template.Must(
	template.New("help").Funcs(
		map[string]any{
			"FormatDefault": formatDefault,
			"FormatRed":     formatRed,
			"FormatBlue":    formatBlue,
			"FormatBold":    formatBold,
		},
	).Parse(helpTemplateRaw),
)
