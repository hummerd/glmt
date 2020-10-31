package impl

import (
	"fmt"
	"net/url"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"
)

// NewMentioner is a factory of mentioners.
func NewMentioner(cfg config.Mentioner) (mentioner.Mentioner, error) {
	const (
		schemeHTTP  = "http"
		schemeHTTPS = "https"
	)

	var dsURL *url.URL
	if cfg.DataSource != "" {
		var err error
		dsURL, err = url.Parse(cfg.DataSource)
		if err != nil {
			return nil, fmt.Errorf("parsing data source: %w", err)
		}
	}

	switch {
	case cfg.MentionsCount < 0:
		return nil, fmt.Errorf("invalid count of mentions: %d", cfg.MentionsCount)
	case !cfg.Enabled, cfg.MentionsCount == 0, dsURL == nil:
		return NewNOOPMentioner(), nil
	case dsURL.Scheme == schemeHTTP, dsURL.Scheme == schemeHTTPS:
		return NewHTTPMentioner(dsURL, cfg.CurrentMemberUID, cfg.MentionsCount), nil
	default:
		return nil, fmt.Errorf("data source %q is not supported", dsURL.Scheme)
	}
}
