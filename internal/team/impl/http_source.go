package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gerr"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

type HTTPSource struct {
	httpClient *http.Client
	url        string
}

func (s *HTTPSource) Team(ctx context.Context) (*team.Team, error) {
	const maxBodySize = 1 << 20

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp *http.Response
	resp, err = s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	defer func() { err = gerr.NewMultiError(err, resp.Body.Close()) }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	var t team.Team
	err = json.NewDecoder(io.LimitReader(resp.Body, maxBodySize)).Decode(&t)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %w", err)
	}

	return &t, nil
}
