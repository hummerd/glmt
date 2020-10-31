package impl

import (
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"
)

func TestMentionDistributor(t *testing.T) {
	const (
		currentMemberUID = "00000000-0000-0000-0000-000000000000"
		currentProject   = "hummerd/glmt"
		testIterations   = 5
	)

	testCases := []mentionDistributorTestCase{{
		Name: "all",

		Members: []mentioner.Member{
			{UID: "00000000-0000-0000-0000-00000000001", IsActive: true},
			{UID: "00000000-0000-0000-0000-00000000002", IsActive: true},
		},
		CurrentMemberUID: currentMemberUID,
		CurrentProject:   currentProject,
		MentionsCount:    2,

		ExpMentionsCount: 2,
		MustIncludeUIDs: []string{
			"00000000-0000-0000-0000-00000000001",
			"00000000-0000-0000-0000-00000000002",
		},
	}, {
		Name: "without current user",

		Members: []mentioner.Member{
			{UID: currentMemberUID, IsActive: true},
			{UID: "00000000-0000-0000-0000-00000000002", IsActive: true},
		},
		CurrentMemberUID: currentMemberUID,
		CurrentProject:   currentProject,
		MentionsCount:    2,

		ExpMentionsCount: 1,
		MustIncludeUIDs: []string{
			"00000000-0000-0000-0000-00000000002",
		},
	}, {
		Name: "with owner",
		Members: []mentioner.Member{
			{UID: "00000000-0000-0000-0000-00000000001", IsActive: true},
			{UID: "00000000-0000-0000-0000-00000000002", IsActive: true},
			{UID: "00000000-0000-0000-0000-00000000004", IsActive: true},
			{UID: "00000000-0000-0000-0000-00000000005", IsActive: true, OwnsProjects: []string{currentProject}},
		},
		CurrentMemberUID: currentMemberUID,
		CurrentProject:   currentProject,
		MentionsCount:    2,

		ExpMentionsCount: 1,
		MustIncludeUIDs:  nil,
	}, {
		Name: "with current user owner",
		Members: []mentioner.Member{
			{UID: currentMemberUID, OwnsProjects: []string{currentProject}},
			{UID: "00000000-0000-0000-0000-00000000002", IsActive: true},
		},
		CurrentMemberUID: currentMemberUID,
		CurrentProject:   currentProject,
		MentionsCount:    2,

		ExpMentionsCount: 1,
		MustIncludeUIDs: []string{
			"00000000-0000-0000-0000-00000000002",
		},
	}, {
		Name: "inactive user",
		Members: []mentioner.Member{
			{UID: "00000000-0000-0000-0000-00000000003", IsActive: false},
			{UID: "00000000-0000-0000-0000-00000000004", IsActive: true},
		},
		CurrentMemberUID: currentMemberUID,
		CurrentProject:   currentProject,
		MentionsCount:    1,

		ExpMentionsCount: 1,
		MustIncludeUIDs: []string{
			"00000000-0000-0000-0000-00000000004",
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			for i := 0; i < testIterations; i++ {
				testMentionDistributor(tt, tc)
			}
		})
	}
}

type mentionDistributorTestCase struct {
	Name string

	Members          []mentioner.Member
	MentionsCount    int
	CurrentMemberUID string
	CurrentProject   string

	ExpMentionsCount int
	MustIncludeUIDs  []string
}

func testMentionDistributor(t *testing.T, tc mentionDistributorTestCase) {
	md := mentionDistributor{
		mentionsConfig: &mentionsConfig{
			Members: tc.Members,
		},
		mentionsCount:    tc.ExpMentionsCount,
		currentMemberUID: tc.CurrentMemberUID,
	}

	gotMembers := md.distribute(tc.CurrentProject)
	if len(gotMembers) != tc.ExpMentionsCount {
		t.Fatalf("mentions count: exp %d, got %d", tc.ExpMentionsCount, len(gotMembers))
	}

	gotMemberUIDs := make(map[string]struct{}, len(gotMembers))
	for _, m := range gotMembers {
		gotMemberUIDs[m.UID] = struct{}{}
	}

	for _, expUID := range tc.MustIncludeUIDs {
		_, ok := gotMemberUIDs[expUID]
		if !ok {
			t.Log(gotMemberUIDs)
			t.Fatalf("exp member: %s", expUID)
		}
	}
}
