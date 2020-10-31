package impl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gerr"
	"gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"
)

func NewHTTPMentioner(url *url.URL, memberUID string, mentionsCount int) *HTTPMentioner {
	return &HTTPMentioner{
		mu: sync.Mutex{},
		mentionDistributor: mentionDistributor{
			mentionsConfig:   nil,
			mentionsCount:    mentionsCount,
			currentMemberUID: memberUID,
		},

		httpClient: http.DefaultClient,
		url:        url,
	}
}

type HTTPMentioner struct {
	// Mu locks mentionDistributor.
	mu sync.Mutex
	mentionDistributor

	httpClient *http.Client
	url        *url.URL
}

// fetchMentionsConfigOnce lazy loads remote config.
func (m *HTTPMentioner) fetchMentionsConfigOnce() (err error) {
	const maxBodySize = 1 << 20

	if m.mentionsConfig != nil {
		return nil
	}

	var resp *http.Response
	resp, err = m.httpClient.Get(m.url.String())
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	defer func() { err = gerr.NewMultiError(err, resp.Body.Close()) }()

	var mentionsConfig mentionsConfig
	err = json.NewDecoder(io.LimitReader(resp.Body, maxBodySize)).Decode(&mentionsConfig)
	if err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	m.mentionsConfig = &mentionsConfig

	return nil
}

// Mentions returns N random members including project owners and
// excluding current member. It is safe for concurrent calls.
func (m *HTTPMentioner) Mentions(project string) (members []mentioner.Member, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	err = m.fetchMentionsConfigOnce()
	if err != nil {
		return nil, fmt.Errorf("fetching remote mentions cfg: %w", err)
	}

	return m.distribute(project), nil
}
