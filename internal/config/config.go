// Package config defines configuration scheme
package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Config struct {
	GitLab    GitLab    `json:"gitlab"`
	MR        MR        `json:"mr"`
	Notifier  Notifier  `json:"notifier"`
	Mentioner Mentioner `json:"mentioner"`
	Hooks     Hooks     `json:"hooks"`
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
	URL     string `json:"url"`
	Message string `json:"message"`
	User    string `json:"user"`
}

type Mentioner struct {
	TeamFileSource string `json:"team_file_source"`
	MentionsCount  int    `json:"count"`
}

type Hooks struct {
	AfterCommands  map[string][]string `json:"after"`
	BeforeCommands map[string][]string `json:"before"`
	Timeout        Duration            `json:"timeout"`
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

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var rawTime string

	err = json.Unmarshal(b, &rawTime)
	if err != nil {
		return err
	}

	var td time.Duration
	td, err = time.ParseDuration(rawTime)
	if err != nil {
		return err
	}

	*d = Duration(td)

	return nil
}
