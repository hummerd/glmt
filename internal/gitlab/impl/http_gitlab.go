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
	"gitlab.com/gitlab-merge-tool/glmt/internal/gitlab"
)

func NewHTTPGitLab(token, host string) *HTTPGitLab {
	return &HTTPGitLab{
		c: &http.Client{
			Timeout: time.Second * 30,
		},
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

	hReq, err := createHTTPRequest(ctx, gl.token, gl.host, req)
	if err != nil {
		return resp, err
	}

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

	resp.URL = createMRURL(gl.host, req.Project, resp.IID)

	return resp, nil
}

func createHTTPRequest(ctx context.Context, token, host string, req gitlab.CreateMRRequest) (*http.Request, error) {
	data := &bytes.Buffer{}
	enc := json.NewEncoder(data)
	enc.SetIndent("", "  ")
	err := enc.Encode(req)
	if err != nil {
		return nil, fmt.Errorf("can not encode request to gitlab: %w", err)
	}

	methodURL := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", host, url.PathEscape(req.Project))
	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, methodURL, data)
	if err != nil {
		return nil, fmt.Errorf("can not create request for gitlab's create mr: %w", err)
	}

	hReq.Header.Set("Content-Type", "application/json")
	hReq.Header.Set("Private-Token", token)

	return hReq, nil
}

func createMRURL(host, project string, iid int64) string {
	return fmt.Sprintf("%s/%s/-/merge_requests/%d", host, project, iid)
}
