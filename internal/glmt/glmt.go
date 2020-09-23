// Package glmt defines logic for glmt tool
package glmt

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/rs/zerolog/log"
	"gitlab.com/glmt/glmt/internal/git"
	"gitlab.com/glmt/glmt/internal/gitlab"
)

func NewGLMT(
	git git.Git,
	gitLab gitlab.GitLab,
) *Core {
	return &Core{
		git:    git,
		gitLab: gitLab,
	}
}

type Core struct {
	git    git.Git
	gitLab gitlab.GitLab
}

type CreateMRParams struct {
	TargetBranch        string
	BranchRegexp        *regexp.Regexp
	TitleTemplate       string
	DescriptionTemplate string
	Squash              bool
	RemoveBranch        bool
}

func (c *Core) CreateMR(ctx context.Context, params *CreateMRParams) error {
	if params.TargetBranch == "" {
		return errors.New("target branch is required")
	}

	br, err := c.git.CurrentBranch()
	if err != nil {
		return err
	}

	r, err := c.git.Remote()
	if err != nil {
		return err
	}

	p, err := projectFromRemote(r)
	if err != nil {
		return err
	}

	ta := getTextArgs(br, p, params)

	t := createText("title", params.TitleTemplate, ta)
	d := createText("description", params.DescriptionTemplate, ta)

	log.Ctx(ctx).Debug().
		Interface("context", ta).
		Str("title", t).
		Str("description", d).
		Msg("create mr")

	_, err = c.gitLab.CreateMR(ctx, gitlab.CreateMRRequest{
		ID:                 p,
		SourceBranch:       br,
		TargetBranch:       params.TargetBranch,
		Title:              t,
		Description:        d,
		Squash:             params.Squash,
		RemoveSourceBranch: params.RemoveBranch,
	})
	if err != nil {
		return err
	}

	return nil
}

func getTextArgs(branch, projectName string, params *CreateMRParams) TextArgs {
	ta := TextArgs{
		ProjectName:      projectName,
		BranchName:       branch,
		TargetBranchName: params.TargetBranch,
	}

	if params.BranchRegexp == nil {
		return ta
	}

	subNames := params.BranchRegexp.SubexpNames()
	if len(subNames) <= 1 {
		return ta
	}

	match := params.BranchRegexp.FindStringSubmatch(branch)
	if len(match) == 0 {
		return ta
	}

	for i := 1; i < len(subNames); i++ {
		switch subNames[i] {
		case "BranchDescription":
			ta.BranchDescription = match[i]

		case "Task":
			ta.Task = match[i]

		case "TaskType":
			ta.TaskType = match[i]
		}
	}

	return ta
}

type TextArgs struct {
	BranchName        string
	TargetBranchName  string
	BranchDescription string
	Task              string
	TaskType          string
	ProjectName       string
}

func createText(part, format string, args TextArgs) string {
	if format == "" {
		return args.BranchName
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

func projectFromRemote(rem string) (string, error) {
	if matchesScheme(rem) {
		url, err := url.Parse(rem)
		if err != nil {
			return "", err
		}

		return url.Path, nil
	}

	if matchesScpLike(rem) {
		p, err := findScpLikePath(rem)
		if err != nil {
			return "", err
		}

		return p, nil
	}

	return "", errors.New("unknown remote path in git repo: " + rem)
}

var (
	isSchemeRegExp   = regexp.MustCompile(`^[^:]+://`)
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})(?:\/|:))?(?P<path>[^\\].*\/[^\\].*)\.git$`)
)

/// findScpLikePath returns the path of the given SCP-like URL.
func findScpLikePath(url string) (string, error) {
	m := scpLikeURLRegExp.FindStringSubmatch(url)

	if len(m) < 5 {
		return "", errors.New("can not find project in remote path: " + url)
	}

	return m[4], nil
}

// MatchesScheme returns true if the given string matches a URL-like
// format scheme.
func matchesScheme(url string) bool {
	return isSchemeRegExp.MatchString(url)
}

// MatchesScpLike returns true if the given string matches an SCP-like
// format scheme.
func matchesScpLike(url string) bool {
	return scpLikeURLRegExp.MatchString(url)
}
