package glmt

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gitlab"
	teami "gitlab.com/gitlab-merge-tool/glmt/internal/team/impl"
)

func TestRemoteParse(t *testing.T) {
	r := "https://github.com/hummerd/client_golang.git"
	p, err := projectFromRemote(r)
	if err != nil {
		t.Fatalf("failed to parse remote %s: %v", r, err)
	}

	if p != "hummerd/client_golang" {
		t.Fatalf("wrong project: %s", p)
	}

	r = "git@bitbucket.org:hummerd/client_golang.git"
	p, err = projectFromRemote(r)
	if err != nil {
		t.Fatalf("failed to parse remote %s: %v", r, err)
	}

	if p != "hummerd/client_golang" {
		t.Fatalf("wrong project: %s", p)
	}

	r = "git@bitbucket.org:hummerd/client_golang"
	p, err = projectFromRemote(r)
	if err != nil {
		t.Fatalf("failed to parse remote %s: %v", r, err)
	}

	if p != "hummerd/client_golang" {
		t.Fatalf("wrong project: %s", p)
	}
}

func TestTextArgs(t *testing.T) {
	expTa := map[string]string{
		"ProjectName":        "prj1",
		"BranchName":         "feature/TASK-123/some-description",
		"TargetBranchName":   "develop",
		"Task":               "TASK-123",
		"TaskType":           "feature",
		"BranchDescription":  "some-description",
		TmpVarGitlabMentions: "",
	}

	params := CreateMRParams{
		TargetBranch: expTa["TargetBranchName"],
		BranchRegexp: regexp.MustCompile(`(?P<TaskType>.*)/(?P<Task>.*)/(?P<BranchDescription>.*)`),
	}

	ta := getTextArgs(expTa["BranchName"], expTa["ProjectName"], params, nil)

	if !reflect.DeepEqual(expTa, ta) {
		t.Fatalf("expected ta: %+v, got %+v", expTa, ta)
	}
}

func TestCreateMR(t *testing.T) {
	gs := &gitStub{
		r: "https://github.com/hummerd/client_golang.git",
		b: "feature/TASK-123/add-some-feature",
	}

	cp := CreateMRParams{
		DescriptionTemplate: "Merge feature {{.Task}} \"{{humanizeText .BranchDescription}}\" into {{.TargetBranchName}}",
		TitleTemplate:       "{{.Task}} {{humanizeText .BranchDescription}}",
		RemoveBranch:        true,
		Squash:              true,
		TargetBranch:        "develop",
		BranchRegexp:        regexp.MustCompile("(?P<TaskType>.*)/(?P<Task>.*)/(?P<BranchDescription>.*)"),
	}

	good := false

	gls := &gitlabStub{
		f: func(method string, arg interface{}) {
			if method == "CreateMR" {
				good = true

				exp := gitlab.CreateMRRequest{
					Description:        "Merge feature TASK-123 \"Add some feature\" into develop",
					Title:              "TASK-123 Add some feature",
					Project:            "hummerd/client_golang",
					SourceBranch:       gs.b,
					TargetBranch:       cp.TargetBranch,
					RemoveSourceBranch: cp.RemoveBranch,
					Squash:             cp.Squash,
				}

				if !reflect.DeepEqual(exp, arg) {
					t.Fatalf("expected create request: %+v, got %+v", exp, arg)
				}
			}
		},
	}

	ts, _ := teami.NewTeamSource("")
	c := Core{
		git:        gs,
		gitLab:     gls,
		teamSource: ts,
	}

	mr, err := c.CreateMR(context.Background(), cp)
	if err != nil {
		t.Fatal("error creating MR", err)
	}

	if mr.ID != 123 {
		t.Fatal("wrong MR id", mr.ID)
	}

	if !good {
		t.Fatal("no call to CreateMR")
	}
}

type gitStub struct {
	r string
	b string
}

func (gs *gitStub) Remote() (string, error) {
	return gs.r, nil
}

func (gs *gitStub) CurrentBranch() (string, error) {
	return gs.b, nil
}

type gitlabCallback func(string, interface{})

type gitlabStub struct {
	f gitlabCallback
}

func (gls *gitlabStub) CreateMR(ctx context.Context, req gitlab.CreateMRRequest) (gitlab.CreateMRResponse, error) {
	gls.f("CreateMR", req)
	return gitlab.CreateMRResponse{
		ID: 123,
	}, nil
}

func (gls *gitlabStub) CurrentUser(ctx context.Context) (gitlab.UserResponse, error) {
	gls.f("CurrentUser", nil)
	return gitlab.UserResponse{
		ID: 123,
	}, nil
}
