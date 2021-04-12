package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/git"
	"gitlab.com/gitlab-merge-tool/glmt/internal/gitlab"
	gitlabi "gitlab.com/gitlab-merge-tool/glmt/internal/gitlab/impl"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	hooksi "gitlab.com/gitlab-merge-tool/glmt/internal/hooks/impl"
	"gitlab.com/gitlab-merge-tool/glmt/internal/notifier"
	notifieri "gitlab.com/gitlab-merge-tool/glmt/internal/notifier/impl"
	teami "gitlab.com/gitlab-merge-tool/glmt/internal/team/impl"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	out := os.Stdout

	var rootCmd = &cobra.Command{Use: "glmt"}
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config")
	rootCmd.PersistentFlags().StringP("token", "k", "", "gitlab API token (get it on /profile/personal_access_tokens page)")
	rootCmd.PersistentFlags().StringP("host", "a", "", "gitlab host")
	rootCmd.PersistentFlags().BoolP("dryrun", "y", false, "dry run true only shows request to gitlab, but do not sends them")
	rootCmd.PersistentFlags().StringP("log", "l", "info", "log level")
	rootCmd.PersistentFlags().Bool("no_hooks", false, "do not run hooks")

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create merge request",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			createMR(cmd, logger, out)
		},
	}
	createFlags := cmdCreate.Flags()
	createFlags.StringP("target", "b", "master", "Merge Request's target branch")
	createFlags.StringP("title", "t", "", "Merge Request's title (template variables can be used in title)")
	createFlags.StringP("description", "d", "", "Merge Request's description (template variables can be used in description)")
	createFlags.StringP("notification_message", "n", "", "Additional notification message")
	rootCmd.AddCommand(cmdCreate)

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Show GLMT version",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			showVersion(cmd, logger, out)
		},
	}
	rootCmd.AddCommand(cmdVersion)

	_ = rootCmd.Execute()
}

func parseLogLevel(flags *pflag.FlagSet) (zerolog.Level, error) {
	log, err := flags.GetString("log")
	if err != nil {
		return zerolog.NoLevel, err
	}

	return zerolog.ParseLevel(log)
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
			return nil, fmt.Errorf("can not read config: %w", err)
		}
	}

	err = applyFlags(flags, cfg)
	if err != nil {
		return nil, fmt.Errorf("can not parse flags: %w", err)
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

func createCore(dryRun bool, out io.StringWriter, cfg *config.Config) (*glmt.Core, error) {
	git, err := git.NewLocalGit()
	if err != nil {
		return nil, err
	}

	gitCfg := cfg.GitLab
	var gitlab gitlab.GitLab
	if dryRun {
		gitlab = gitlabi.NewDryRunGitLab(out, gitCfg.Token, gitCfg.URL)
	} else {
		gitlab = gitlabi.NewHTTPGitLab(gitCfg.Token, gitCfg.URL)
	}

	nfyCfg := cfg.Notifier
	var ns []notifier.Notifier
	if nfyCfg.SlackWebHook.Enabled {
		ns = append(ns, notifieri.NewSlackWebHookNotifier(nfyCfg.SlackWebHook))
	}
	if nfyCfg.Telegram.Enabled {
		ns = append(ns, notifieri.NewTelegramNotifier(nfyCfg.Telegram))
	}
	n := notifieri.NewMultiNotifier(ns...)

	mrCfg := cfg.Mentioner
	ts, err := teami.NewTeamSource(mrCfg.TeamFileSource)
	if err != nil {
		return nil, err
	}

	hsCfg := cfg.Hooks
	hs := hooksi.NewHooks(hsCfg, os.Stdout, os.Stderr)

	return glmt.NewGLMT(git, gitlab, n, ts, hs), nil
}
