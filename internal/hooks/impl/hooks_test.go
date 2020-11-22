package impl_test

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/hooks"
	"gitlab.com/gitlab-merge-tool/glmt/internal/hooks/impl"
)

func TestHooks_RunAfter(t *testing.T) {
	h := impl.NewHooks(config.Hooks{
		AfterCommands: map[string][]string{
			"test": []string{"sh"},
		},
	}, nil, nil)

	err := h.RunAfter(context.Background(), hooks.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHooks_RunBefore_Output(t *testing.T) {
	const expBuf = "Hello World\n"

	var buf bytes.Buffer
	h := impl.NewHooks(config.Hooks{
		BeforeCommands: map[string][]string{
			"test": []string{
				"echo", "Hello", "World",
			},
		},
	}, &buf, &buf)

	err := h.RunBefore(context.Background(), hooks.Params{})
	switch {
	case err != nil:
		t.Fatal(err)
	case buf.String() != expBuf:
		t.Fatal("Invalid output", buf.String())
	}
}

func TestHooks_RunBefore_Fail(t *testing.T) {
	h := impl.NewHooks(config.Hooks{
		BeforeCommands: map[string][]string{
			"test": []string{"this-command-should-not-exists"},
		},
	}, nil, nil)

	err := h.RunBefore(context.Background(), hooks.Params{})
	if !errors.Is(err, exec.ErrNotFound) {
		t.Fatal(err)
	}
}

func TestHooks_RunBefore_Env(t *testing.T) {
	const expBuf = "master\n"

	var buf bytes.Buffer
	h := impl.NewHooks(config.Hooks{
		BeforeCommands: map[string][]string{
			"test": []string{
				"sh", "-c", "echo $GLMT_BRANCH",
			},
		},
	}, &buf, &buf)

	err := h.RunBefore(context.Background(), hooks.Params{
		Branch: "master",
	})
	switch {
	case err != nil:
		t.Log(buf.String())
		t.Fatal(err)
	case buf.String() != expBuf:
		t.Fatal("Invalid output", buf.String())
	}
}

func TestHooks_RunAfter_Deadline(t *testing.T) {
	var buf bytes.Buffer
	h := impl.NewHooks(config.Hooks{
		AfterCommands: map[string][]string{
			"test": []string{
				"sleep", "1",
			},
		},
		Timeout: config.Duration(100 * time.Millisecond),
	}, &buf, &buf)

	err := h.RunAfter(context.Background(), hooks.Params{})
	if err == nil {
		t.Fatal("No error")
	}
}
