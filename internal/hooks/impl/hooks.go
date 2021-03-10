// Package impl implements hooks.Runner.
package impl

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/hooks"
)

type Hooks struct {
	afterCommands  map[string][]string
	beforeCommands map[string][]string

	timeout time.Duration

	stdout io.Writer
	stderr io.Writer
}

func NewHooks(cfg config.Hooks, stdout io.Writer, stderr io.Writer) *Hooks {
	const defaultTimeout = 5 * time.Minute

	timeout := time.Duration(cfg.Timeout)
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &Hooks{
		afterCommands:  filterCommands(cfg.AfterCommands),
		beforeCommands: filterCommands(cfg.BeforeCommands),

		timeout: timeout,

		stdout: stdout,
		stderr: stderr,
	}
}

func (h Hooks) RunAfter(ctx context.Context, params hooks.Params) (err error) {
	return h.run(ctx, h.afterCommands, params)
}

func (h Hooks) RunBefore(ctx context.Context, params hooks.Params) error {
	return h.run(ctx, h.beforeCommands, params)
}

func (h Hooks) run(
	ctx context.Context,
	commands map[string][]string,
	params hooks.Params,
) (err error) {
	env := params.Env()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, h.timeout)
	defer cancel()

	for name, cmd := range commands {
		cmdProc := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
		cmdProc.Env = append(os.Environ(), env...)
		cmdProc.Stdout = h.stdout
		cmdProc.Stderr = h.stderr

		err = cmdProc.Run()
		if err != nil {
			return fmt.Errorf("%s: running command: %w", name, err)
		}
	}

	return nil
}

// filterCommands removes empty commands.
func filterCommands(commands map[string][]string) (filtered map[string][]string) {
	filtered = make(map[string][]string, len(commands))

	for name, cmd := range commands {
		if len(cmd) == 0 {
			continue
		}

		filtered[name] = cmd
	}

	return filtered
}
