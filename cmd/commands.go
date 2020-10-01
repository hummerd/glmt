package main

import (
	"context"
	"io"
	"os"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"gitlab.com/glmt/glmt/internal/glmt"
)

func createMR(cmd *cobra.Command, args []string, logger zerolog.Logger, out io.StringWriter) {
	cfg, err := finalConfig(cmd.Flags())
	if err != nil {
		panic(err)
	}

	ctx := logger.WithContext(context.Background())

	core := createCore(cfg)

	br, err := regexp.Compile(cfg.MR.BranchRegexp)
	if err != nil {
		_, _ = out.WriteString("Failed to compile branch regexp: " + err.Error() + "\n")
		os.Exit(1)
	}

	params := &glmt.CreateMRParams{
		TargetBranch:        cfg.MR.TargetBranch,
		BranchRegexp:        br,
		TitleTemplate:       cfg.MR.Title,
		DescriptionTemplate: cfg.MR.Description,
		Squash:              cfg.MR.Squash,
		RemoveBranch:        cfg.MR.RemoveSourceBranch,
	}

	mr, err := core.CreateMR(ctx, params)
	if err != nil {
		_, _ = out.WriteString("Failed to create MR: " + err.Error() + "\n")
		os.Exit(1)
	}

	_, _ = out.WriteString("MR created\n")
	_, _ = out.WriteString(mr.URL + "\n")
}
