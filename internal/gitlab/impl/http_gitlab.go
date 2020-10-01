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

func (gl *HTTPGitLab) CreateMR(ctx context.Context, req gitlab.CreateMRRequest) (gitlab.CreateMRResponse, error) {
	var resp gitlab.CreateMRResponse

	data := &bytes.Buffer{}
	err := json.NewEncoder(data).Encode(req)
	if err != nil {
		return resp, fmt.Errorf("can not encode request to gitlab: %w", err)
	}

	methodURL := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", gl.host, url.PathEscape(req.Project))
	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, methodURL, data)
	if err != nil {
		return resp, fmt.Errorf("can not create request for gitlab's create mr: %w", err)
	}

	hReq.Header.Set("Content-Type", "application/json")
	hReq.Header.Set("Private-Token", gl.token)

	buff := &bytes.Buffer{}
	err = json.NewEncoder(buff).Encode(req)
	if err != nil {
		return resp, fmt.Errorf("can not create request payload for gitlab's create mr: %w", err)
	}

	hReq.Body = ioutil.NopCloser(buff)

	log.Ctx(ctx).Debug().
		Stringer("url", hReq.URL).
		Str("token", hideToken(gl.token)).
		Interface("request", req).
		Msg("post to gitlab")

	hResp, err := gl.c.Do(hReq)
	if err != nil {
		return resp, fmt.Errorf("can not create mr in gitlab: %w", err)
	}

	defer hResp.Body.Close()

	if hResp.StatusCode != http.StatusCreated {
		// TODO: implement typed error
		// type glerr  struct {
		// 	Message map[string]interface{}
		// 	Error string
		// }
		// err = json.NewDecoder(hResp.Body).Decode(&gerr)

		errm, err := ioutil.ReadAll(hResp.Body)
		if err != nil {
			return resp, fmt.Errorf("can not decode error from gitlab's create MR: %w", err)
		}
		return resp, gitlab.GitlabError{Message: string(errm)}
	}

	err = json.NewDecoder(hResp.Body).Decode(&resp)
	if err != nil {
		return resp, fmt.Errorf("can not decode response from gitlab's create MR: %w", err)
	}

	resp.URL = fmt.Sprintf("%s/%s/-/merge_requests/%d", gl.host, req.Project, resp.IID)

	return resp, nil
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
