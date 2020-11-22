package impl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"
)

const (
	// %s - is a bot api key.
	telegramSendMessagePath = "bot%s/sendMessage"
)

func NewTelegramNotifier(cfg config.Telegram) *TelegramNotifier {
	const defaultMessageTmpl = "{{.Description}}\n{{.MergeRequestURL}}"

	if cfg.MessageTmpl == "" {
		cfg.MessageTmpl = defaultMessageTmpl
	}

	return &TelegramNotifier{
		url:         cfg.URL,
		apiKey:      cfg.APIKey,
		messageTmpl: cfg.MessageTmpl,
		chatID:      cfg.ChatID,

		httpClient: http.DefaultClient,
	}
}

type TelegramNotifier struct {
	url         string
	apiKey      string
	chatID      string
	messageTmpl string

	httpClient *http.Client
}

func (tn *TelegramNotifier) Send(
	ctx context.Context,
	args map[string]string,
	add string,
	mentions []*team.Member,
) error {
	args[glmt.TmpVarNotificationMentions] = getMentions(
		mentions,
		memberKeyTelegram,
		"@%s",
	)

	m := templating.CreateText("telegram_wh_message", tn.messageTmpl, args)

	if add != "" {
		m += "\n" + add
	}

	u, err := url.Parse(tn.url)
	if err != nil {
		return fmt.Errorf("parsing api url: %w", err)
	}

	u.Path = fmt.Sprintf(telegramSendMessagePath, tn.apiKey)

	query := u.Query()
	query.Set("chat_id", tn.chatID)
	query.Set("text", m)
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("doing get request: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("closing body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("wrong status: %s", resp.Status)
	}

	return nil
}
