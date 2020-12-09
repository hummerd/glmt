package impl

import (
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

func TestGetMentions(t *testing.T) {
	const exp = "<@test_id1>, <@test_id2>"
	const format = "<@%s>"
	const memberKey = "test"

	got := getMentions([]*team.Member{{
		Names: map[string]string{memberKey: "test_id1"},
	}, {
		Names: map[string]string{memberKey: "test_id2"},
	}}, memberKey, format)

	if exp != got {
		t.Fatalf("exp: %s, got: %s", exp, got)
	}
}
