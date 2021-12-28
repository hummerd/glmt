// Package config defines configuration scheme
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Config struct {
	Base      string    `json:"base"`
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
	BranchRegexp       string   `json:"branch_regexp"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	TargetBranch       string   `json:"target_branch"`
	Squash             bool     `json:"squash"`
	RemoveSourceBranch bool     `json:"remove_source_branch"`
	LabelVars          []string `json:"label_vars"`
}

type Notifier struct {
	SlackWebHook      SlackWebHook      `json:"slack_web_hook"`
	Telegram          Telegram          `json:"telegram"`
	MattermostWebHook MattermostWebHook `json:"mattermost_web_hook"`
}

type MattermostWebHook struct {
	Enabled     bool   `json:"enabled"`
	URL         string `json:"url"`
	MessageTmpl string `json:"message"`
	User        string `json:"user"`
}

type SlackWebHook struct {
	Enabled     bool   `json:"enabled"`
	URL         string `json:"url"`
	MessageTmpl string `json:"message"`
	User        string `json:"user"`
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

type Telegram struct {
	Enabled     bool   `json:"enabled"`
	URL         string `json:"url"`
	APIKey      string `json:"api_key"`
	MessageTmpl string `json:"message"`
	ChatID      string `json:"chat_id"`
}

func LoadConfig(path string) (*Config, error) {
	c, err := loadFromFile(path)
	if err != nil {
		return nil, err
	}

	if c.Base != "" {
		var bc *Config
		if strings.HasPrefix(c.Base, "http://") ||
			strings.HasPrefix(c.Base, "https://") {
			bc, err = loadFromHttp(c.Base)
		} else {
			bc, err = loadFromFile(c.Base)
		}

		if err != nil {
			return nil, err
		}

		err = mergo.Merge(c, bc)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config: %w", err)
		}
	}

	return c, nil
}

func loadFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	defer f.Close()

	var c Config
	err = json5.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}
	return &c, nil
}

func loadFromHttp(url string) (*Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad response from config source: " + resp.Status)
	}

	var c Config
	err = json5.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
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
