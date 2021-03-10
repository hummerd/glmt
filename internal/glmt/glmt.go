// Package glmt defines logic for glmt tool
package glmt

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gerr"
	"gitlab.com/gitlab-merge-tool/glmt/internal/gitlab"
	"gitlab.com/gitlab-merge-tool/glmt/internal/hooks"
	"gitlab.com/gitlab-merge-tool/glmt/internal/notifier"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"
)

var (
	Version     string
	VersionDate string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ErrNotification = errors.New("notification error")
)

func NewGLMT(
	git Git,
	gitLab gitlab.GitLab,
	notifier notifier.Notifier,
	teamSource team.TeamFileSource,
	hooks hooks.Runner,
) *Core {
	return &Core{
		git:        git,
		gitLab:     gitLab,
		notifier:   notifier,
		teamSource: teamSource,
		hooks:      hooks,
	}
}

type Core struct {
	git        Git
	gitLab     gitlab.GitLab
	notifier   notifier.Notifier
	teamSource team.TeamFileSource
	hooks      hooks.Runner
}

type CreateMRParams struct {
	TargetBranch        string
	BranchRegexp        *regexp.Regexp
	TitleTemplate       string
	DescriptionTemplate string
	Squash              bool
	RemoveBranch        bool
	NotificationMessage string
	MentionsCount       int
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

	cu, err := c.gitLab.CurrentUser(ctx)
	if err != nil {
		return mr, err
	}

	var ms []*team.Member
	if c.teamSource != nil && params.MentionsCount > 0 {
		tm, err := c.teamSource.Team(ctx)
		if err != nil {
			return mr, err
		}

		ms = Mentions(tm, cu.Username, p, params.MentionsCount)
	}

	ta := getTextArgs(br, p, r, cu.Username, params, ms)

	var t string
	if params.TitleTemplate != "" {
		t = templating.CreateText("title", params.TitleTemplate, ta)
	}

	t = strings.TrimSpace(t)
	if t == "" {
		t = br
	}

	var d string
	if params.DescriptionTemplate != "" {
		d = templating.CreateText("description", params.DescriptionTemplate, ta)
	}

	d = strings.TrimSpace(d)
	if d == "" {
		d = "Merge " + br + " into " + params.TargetBranch
	}

	err = c.hooks.RunBefore(ctx, hooks.Params(ta))
	if err != nil {
		return mr, fmt.Errorf("hooks precondition failed: %w", err)
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
		AssigneeID:         cu.ID,
	})
	if err != nil {
		return mr, err
	}

	ta[TmpVarTitle] = t
	ta[TmpVarDescription] = d
	ta[TmpVarMRURL] = gmr.URL
	ta[TmpVarMRChangesCount] = gmr.ChangesCount

	err = c.hooks.RunAfter(ctx, hooks.Params(ta))
	if err != nil {
		return mr, fmt.Errorf("hooks postcondition failed: %w", err)
	}

	mr.ID = gmr.ID
	mr.IID = gmr.IID
	mr.ProjectID = gmr.ProjectID
	mr.CreatedAt = gmr.CreatedAt
	mr.URL = gmr.URL

	if c.notifier != nil {
		err = c.notifier.Send(ctx, ta, params.NotificationMessage, ms)
		if err != nil {
			err = gerr.NewNestedError(ErrNotification, err)
		}

		log.Ctx(ctx).Debug().
			Interface("context", ta).
			Msg("notification")
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
