package hooks_test

import (
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/hooks"
)

func TestHookParamsEnv(t *testing.T) {
	const exp = "GLMT_BRANCH=master"

	params := hooks.Params{
		"Branch": "master",
	}

	env := params.Env()

	var got bool
	for _, envVal := range env {
		if envVal == exp {
			got = true
			break
		}
	}

	if !got {
		t.Fatal(env)
	}
}
