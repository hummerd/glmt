package glmt

import (
	"bytes"
	"strings"
	"text/template"
	"unicode"
)

const (
	TmpVarProjectName      = "ProjectName"
	TmpVarBranchName       = "BranchName"
	TmpVarTargetBranchName = "TargetBranchName"
)

func getTextArgs(branch, projectName string, params CreateMRParams) map[string]string {
	r := map[string]string{}

	defer func() {
		// in the end override values with well known
		r[TmpVarProjectName] = projectName
		r[TmpVarBranchName] = branch
		r[TmpVarTargetBranchName] = params.TargetBranch
	}()

	if params.BranchRegexp == nil {
		return r
	}

	subNames := params.BranchRegexp.SubexpNames()
	if len(subNames) <= 1 {
		return r
	}

	match := params.BranchRegexp.FindStringSubmatch(branch)
	if len(match) == 0 {
		return r
	}

	for i := 1; i < len(subNames); i++ {
		r[subNames[i]] = match[i]
	}

	return r
}

func createText(part, format string, args map[string]string) string {
	if format == "" {
		return args[TmpVarBranchName]
	}

	funcMap := template.FuncMap{
		"humanizeText": humanizeText,
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
