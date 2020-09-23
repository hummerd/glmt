package glmt

import (
	"regexp"
	"testing"
)

func TestTextArgs(t *testing.T) {
	expTa := TextArgs{
		ProjectName:       "prj1",
		BranchName:        "feature/TASK-123/some-description",
		TargetBranchName:  "develop",
		Task:              "TASK-123",
		TaskType:          "feature",
		BranchDescription: "some-description",
	}

	params := &CreateMRParams{
		TargetBranch: expTa.TargetBranchName,
		BranchRegexp: regexp.MustCompile(`(?P<TaskType>.*)/(?P<Task>.*)/(?P<BranchDescription>.*)`),
	}

	ta := getTextArgs(expTa.BranchName, expTa.ProjectName, params)

	if ta != expTa {
		t.Fatalf("expected ta: %+v, got %+v", expTa, ta)
	}
}
