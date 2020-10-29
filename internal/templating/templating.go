package templating

import (
	"bytes"
	"html/template"
	"strings"
	"unicode"
)

func CreateText(part, format string, args map[string]string) string {
	if format == "" {
		return ""
	}

	funcMap := template.FuncMap{
		"humanizeText": humanizeText,
		"upper":        strings.ToUpper,
		"lower":        strings.ToLower,
	}

	tmpl, _ := template.New(part).Funcs(funcMap).Parse(format)

	buff := &bytes.Buffer{}
	_ = tmpl.Execute(buff, args)

	return buff.String()
}

func isSeparator(r rune) bool {
	switch {
	case r == '_':
		return true
	case r == '-':
		return true
	}

	return false
}

func humanizeText(s string) string {
	first := true
	return strings.Map(
		func(r rune) rune {
			if isSeparator(r) {
				return ' '
			}

			if first && unicode.IsLetter(r) {
				first = false
				return unicode.ToTitle(r)
			}
			return r
		},
		s)
}
