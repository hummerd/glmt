// Package glmt defines logic for glmt tool
package glmt

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/glmt/glmt/internal/gitlab"
	"gitlab.com/glmt/glmt/internal/templating"
)

var (
	ErrNotification = errors.New("notification error")
)

func NewGLMT(
	git Git,
	gitLab gitlab.GitLab,
	notifier Notifier,
) *Core {
	return &Core{
		git:      git,
		gitLab:   gitLab,
		notifier: notifier,
	}
}

type Core struct {
	git      Git
	gitLab   gitlab.GitLab
	notifier Notifier
}

type CreateMRParams struct {
	TargetBranch        string
	BranchRegexp        *regexp.Regexp
	TitleTemplate       string
	DescriptionTemplate string
	Squash              bool
	RemoveBranch        bool
	NotificationMessage string
}

type MergeRequest struct {
	ID        int64     `json:"id"`
	IID       int64     `json:"iid"`
	ProjectID int64     `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

func (c *Core) CreateMR(ctx context.Context, params CreateMRParams) (MergeRequest, error) {
	var mr MergeRequest
	if params.TargetBranch == "" {
		return mr, errors.New("target branch is required")
	}

	br, err := c.git.CurrentBranch()
	if err != nil {
		return mr, err
	}

	r, err := c.git.Remote()
	if err != nil {
		return mr, err
	}

	p, err := projectFromRemote(r)
	if err != nil {
		return mr, err
	}

	ta := getTextArgs(br, p, params)

	var t string
	if params.TitleTemplate != "" {
		t = templating.CreateText("title", params.TitleTemplate, ta)
	}

	var d string
	if params.DescriptionTemplate != "" {
		d = templating.CreateText("description", params.DescriptionTemplate, ta)
	} else {
		d = "Merge " + br + " into " + params.TargetBranch
	}

	log.Ctx(ctx).Debug().
		Interface("context", ta).
		Str("title", t).
		Str("description", d).
		Msg("create mr")

	gmr, err := c.gitLab.CreateMR(ctx, gitlab.CreateMRRequest{
		Project:            p,
		SourceBranch:       br,
		TargetBranch:       params.TargetBranch,
		Title:              t,
		Description:        d,
		Squash:             params.Squash,
		RemoveSourceBranch: params.RemoveBranch,
	})
	if err != nil {
		return mr, err
	}

	mr.ID = gmr.ID
	mr.IID = gmr.IID
	mr.ProjectID = gmr.ProjectID
	mr.CreatedAt = gmr.CreatedAt
	mr.URL = gmr.URL

	if c.notifier != nil {
		ta[TmpVarTitle] = t
		ta[TmpVarDescription] = d
		ta[TmpVarMRURL] = gmr.URL

		err = c.notifier.Send(ctx, ta, params.NotificationMessage)
		if err != nil {
			err = NewNestedError(ErrNotification, err)
		}
	}

	return mr, err
}

func projectFromRemote(rem string) (string, error) {
	var p string
	if matchesScheme(rem) {
		url, err := url.Parse(rem)
		if err != nil {
			return "", err
		}

		p = strings.TrimLeft(url.Path, "/")
	} else if matchesScpLike(rem) {
		var err error
		p, err = findScpLikePath(rem)
		if err != nil {
			return "", err
		}
	}

	if p != "" {
		if strings.HasSuffix(p, ".git") {
			p = p[:len(p)-4]
		}
		return p, nil
	}

	return "", errors.New("unknown remote path in git repo: " + rem)
}

var (
	isSchemeRegExp   = regexp.MustCompile(`^[^:]+://`)
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})(?:\/|:))?(?P<path>[^\\].*\/[^\\].*)$`)
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
