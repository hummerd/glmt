// Package impl implements http service for gitlab
package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/glmt/glmt/internal/gitlab"
)

func NewHTTPGitLab(token, host string) *HTTPGitLab {
	return &HTTPGitLab{
		c:     &http.Client{},
		token: token,
		host:  host,
	}
}

type HTTPGitLab struct {
	c     *http.Client
	token string
	host  string
}

func (gl *HTTPGitLab) CreateMR(ctx context.Context, req gitlab.CreateMRRequest) (int64, error) {
	data := &bytes.Buffer{}
	err := json.NewEncoder(data).Encode(req)
	if err != nil {
		return 0, fmt.Errorf("can not encode request to gitlab: %w", err)
	}

	methodURL := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", gl.host, url.PathEscape(req.ID))
	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, methodURL, data)
	if err != nil {
		return 0, fmt.Errorf("can not create request for gitlab's create mr: %w", err)
	}

	hReq.Header.Set("Content-Type", "application/json")
	hReq.Header.Set("Private-Token", gl.token)

	buff := &bytes.Buffer{}
	err = json.NewEncoder(buff).Encode(req)
	if err != nil {
		return 0, fmt.Errorf("can not create request payload for gitlab's create mr: %w", err)
	}

	hReq.Body = ioutil.NopCloser(buff)

	log.Ctx(ctx).Debug().
		Stringer("url", hReq.URL).
		Str("token", hideToken(gl.token)).
		Interface("request", req).
		Msg("post to gitlab")

	hResp, err := gl.c.Do(hReq)
	if err != nil {
		return 0, fmt.Errorf("can not create mr in gitlab: %w", err)
	}

	type tResp struct {
		ID        int64     `json:"id"`
		IID       int64     `json:"iid"`
		ProjectID int64     `json:"project_id"`
		CreatedAt time.Time `json:"created_at"`
	}
	var resp tResp

	br, _ := ioutil.ReadAll(hResp.Body)
	fmt.Println(string(br))

	defer hResp.Body.Close()
	err = json.NewDecoder(hResp.Body).Decode(&resp)
	if err != nil {
		return 0, fmt.Errorf("can not decode response from gitlab's create MR: %w", err)
	}

	return resp.ID, nil
}

func hideToken(s string) string {
	if s == "" {
		return s
	}

	if len(s) < 3 {
		return s
	}

	return s[:3] + "..."
}
