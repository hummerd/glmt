package impl

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gitlab"
)

func NewDryRunGitLab(out io.StringWriter, token, host string) *DryRunGitLab {
	return &DryRunGitLab{
		out:   out,
		token: token,
		host:  host,
	}
}

type DryRunGitLab struct {
	out   io.StringWriter
	token string
	host  string
}

func (gl *DryRunGitLab) CreateMR(ctx context.Context, req gitlab.CreateMRRequest) (gitlab.CreateMRResponse, error) {
	var resp gitlab.CreateMRResponse

	hReq, err := createHTTPRequest(ctx, gl.token, gl.host, req)
	if err != nil {
		return resp, err
	}

	_, _ = gl.out.WriteString("Sending create request:\n")
	writeRequest(gl.out, hReq)

	return gitlab.CreateMRResponse{
		URL: createMRURL(gl.host, req.Project, resp.IID),
	}, nil
}

func createMRURL(host, project string, iid int64) string {
	return fmt.Sprintf("%s/%s/-/merge_requests/%d", host, project, iid)
}

func (gl *DryRunGitLab) CurrentUser(ctx context.Context) (gitlab.UserResponse, error) {
	return gitlab.UserResponse{}, nil
}

func writeRequest(out io.StringWriter, r *http.Request) {
	_, _ = out.WriteString(fmt.Sprintf("%v %v\n", r.Method, r.URL))

	for name, headers := range r.Header {
		for _, h := range headers {
			if name == "Private-Token" {
				h = hideToken(h)
			}

			_, _ = out.WriteString(fmt.Sprintf("%v: %v\n", name, h))
		}
	}

	body, _ := ioutil.ReadAll(r.Body)
	_, _ = out.WriteString(string(body))
	_, _ = out.WriteString("\n")
}
