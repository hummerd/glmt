// Package config defines configuration scheme
package config

import (
	"os"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Config struct {
	GitLab    GitLab    `json:"gitlab"`
	MR        MR        `json:"mr"`
	Notifier  Notifier  `json:"notifier"`
	Mentioner Mentioner `json:"mentioner"`
}

type Mentioner struct {
	Enabled          bool   `json:"enabled"`
	DataSource       string `json:"data_source"`
	MentionsCount    int    `json:"count"`
	CurrentMemberUID string `json:"member_uid"`
}

type GitLab struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type MR struct {
	BranchRegexp       string `json:"branch_regexp"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	TargetBranch       string `json:"target_branch"`
	Squash             bool   `json:"squash"`
	RemoveSourceBranch bool   `json:"remove_source_branch"`
}

type Notifier struct {
	SlackWebHook SlackWebHook `json:"slack_web_hook"`
}

type SlackWebHook struct {
	URL             string `json:"url"`
	MessageTemplate string `json:"message_template"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var c Config
	err = json5.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
