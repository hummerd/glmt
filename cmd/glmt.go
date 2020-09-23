package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/glmt/glmt/internal/config"
	giti "gitlab.com/glmt/glmt/internal/git/impl"
	gitlabi "gitlab.com/glmt/glmt/internal/gitlab/impl"
	"gitlab.com/glmt/glmt/internal/glmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	logger := zerolog.New(os.Stdout)
	out := os.Stdout

	var rootCmd = &cobra.Command{Use: "glmt"}
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config")
	rootCmd.PersistentFlags().StringP("token", "k", "", "gitlab API token")
	rootCmd.PersistentFlags().StringP("host", "a", "", "gitlab host")

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create merge request",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			createMR(cmd, args, logger, out)
		},
	}
	createFlags := cmdCreate.Flags()
	createFlags.StringP("target", "b", "master", "Merge Request's target branch")
	createFlags.StringP("title", "t", "", "Merge Request's title (template variables can be used in title)")
	createFlags.StringP("description", "d", "", "Merge Request's description (template variables can be used in description)")

	rootCmd.AddCommand(cmdCreate)
	_ = rootCmd.Execute()
}

func finalConfig(flags *pflag.FlagSet) (*config.Config, error) {
	cp, err := flags.GetString("config")
	if err != nil {
		return nil, err
	}

	defaultCfg := false
	if cp == "" {
		cd, _ := os.UserConfigDir()
		cp = filepath.Join(cd, "glmt.config")
		defaultCfg = true
	}

	cfg := &config.Config{}

	if _, err := os.Stat(cp); err != nil {
		if os.IsNotExist(err) {
			if !defaultCfg {
				return nil, errors.New("config does not exists in: " + cp)
			}
		} else {
			return nil, fmt.Errorf("can not read config: %s, %w", cp, err)
		}
	} else {
		cfg, err = config.LoadConfig(cp)
		if err != nil {
			return nil, errors.New("can not read config: " + err.Error())
		}
	}

	err = applyFlags(flags, cfg)
	if err != nil {
		return nil, errors.New("can not parse flags: " + err.Error())
	}

	return cfg, nil
}

func applyFlags(flags *pflag.FlagSet, cfg *config.Config) error {
	t, err := flags.GetString("token")
	if err != nil {
		return err
	}

	if t != "" {
		cfg.GitLab.Token = t
	}

	h, err := flags.GetString("host")
	if err != nil {
		return err
	}

	if h != "" {
		cfg.GitLab.URL = t
	}
	if cfg.GitLab.URL == "" {
		cfg.GitLab.URL = "https://gitlab.com"
	}

	mrt, err := flags.GetString("title")
	if err != nil {
		return err
	}

	if mrt != "" {
		cfg.MR.Title = mrt
	}

	mrd, err := flags.GetString("description")
	if err != nil {
		return err
	}

	if mrd != "" {
		cfg.MR.Description = mrd
	}

	var target string
	if flags.Changed("target") {
		target, err = flags.GetString("target")
		if err != nil {
			return err
		}
	}

	if target != "" {
		cfg.MR.TargetBranch = target
	}

	return nil
}

func createCore(cfg *config.Config) *glmt.Core {
	git, _ := giti.NewLocalGit()
	gitlab := gitlabi.NewHTTPGitLab(cfg.GitLab.Token, cfg.GitLab.URL)

	return glmt.NewGLMT(git, gitlab)
}
