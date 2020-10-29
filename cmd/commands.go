package main

import (
	"context"
	"io"
	"os"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
)

func createMR(cmd *cobra.Command, logger zerolog.Logger, out io.StringWriter) {
	flags := cmd.Flags()
	cfg, err := finalConfig(flags)
	if err != nil {
		panic(err)
	}

	ll, err := parseLogLevel(flags)
	if err != nil {
		_, _ = out.WriteString("Failed to parse log level: " + err.Error() + "\n")
		os.Exit(1)
	}

	logger = logger.Level(ll)
	ctx := logger.WithContext(context.Background())

	dryRun, err := flags.GetBool("dryrun")
	if err != nil {
		_, _ = out.WriteString("Failed to parse dryrun: " + err.Error() + "\n")
		os.Exit(1)
	}

	core, err := createCore(dryRun, out, &cfg.GitLab, &cfg.Notifier)
	if err != nil {
		_, _ = out.WriteString("Failed to start glmt: " + err.Error() + "\n")
		os.Exit(1)
	}

	br, err := regexp.Compile(cfg.MR.BranchRegexp)
	if err != nil {
		_, _ = out.WriteString("Failed to compile branch regexp: " + err.Error() + "\n")
		os.Exit(1)
	}

	na, err := flags.GetString("notification_message")
	if err != nil {
		_, _ = out.WriteString("Failed to parse notification_message: " + err.Error() + "\n")
		os.Exit(1)
	}

	params := glmt.CreateMRParams{
		TargetBranch:        cfg.MR.TargetBranch,
		BranchRegexp:        br,
		TitleTemplate:       cfg.MR.Title,
		DescriptionTemplate: cfg.MR.Description,
		Squash:              cfg.MR.Squash,
		RemoveBranch:        cfg.MR.RemoveSourceBranch,
		NotificationMessage: na,
	}

	mr, err := core.CreateMR(ctx, params)
	if err != nil {
		_, _ = out.WriteString("Failed to create MR: " + err.Error() + "\n")
		os.Exit(1)
	}

	_, _ = out.WriteString("MR created\n")
	_, _ = out.WriteString(mr.URL + "\n")
}
