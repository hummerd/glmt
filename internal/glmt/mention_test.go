package glmt

import (
	"reflect"
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

func TestMentioner(t *testing.T) {
	const (
		currentProject = "hummerd/glmt"
	)

	m1 := team.Member{
		Username:     "billi",
		IsActive:     true,
		OwnsProjects: []string{currentProject},
	}
	m2 := team.Member{
		Username: "allan",
		IsActive: false,
	}
	m3 := team.Member{
		Username: "william",
		IsActive: true,
	}
	m4 := team.Member{
		Username:     "robert",
		IsActive:     true,
		OwnsProjects: []string{currentProject},
	}
	expMembers := []*team.Member{&m4, &m3}

	tm := team.Team{
		Members: []*team.Member{&m1, &m2, &m3, &m4},
	}

	ms := Mentions(&tm, "@billi", currentProject, 2)
	if !reflect.DeepEqual(expMembers, ms) {
		t.Fatalf("exp: %v, got: %v", expMembers, ms)
	}
}
