// Package impl is implementation of teamfile package
package impl

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

func NewTeamSource(src string) (team.TeamFileSource, error) {
	const (
		schemeHTTP  = "http"
		schemeHTTPS = "https"
	)

	var dsURL *url.URL
	if src != "" {
		var err error
		dsURL, err = url.Parse(src)
		if err != nil {
			return nil, fmt.Errorf("can not parse team source: %w", err)
		}
	}

	switch {
	case src == "":
		return nil, nil
	case dsURL.Scheme == schemeHTTP, dsURL.Scheme == schemeHTTPS:
		return &HTTPSource{
			httpClient: &http.Client{
				Timeout: time.Second * 5,
			},
			url: src,
		}, nil
	default:
		_, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("team file %q not exists", src)
			} else {
				return nil, fmt.Errorf("team file source %q not supported: %w", src, err)
			}
		}

		return &LocalSource{
			path: src,
		}, nil
	}
}
