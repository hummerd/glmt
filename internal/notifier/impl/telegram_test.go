package impl_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	"gitlab.com/gitlab-merge-tool/glmt/internal/notifier/impl"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

func TestTelegramNotifier_Send(t *testing.T) {
	const (
		description = "test_description"
		addText     = "test_add"
		username    = "test_username"
	)

	cfg := config.Telegram{
		Enabled:     true,
		APIKey:      "TEST_API_KEY",
		ChatID:      "-00000000000",
		MessageTmpl: "{{.Description}}\n{{.NotificationMentions}}",
		// URL will be set after test server initialization.
		URL: "",
	}

	tsDone := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		defer close(tsDone)

		t.Log(r.URL)

		switch {
		case !strings.Contains(r.URL.Path, cfg.APIKey):
			t.Fatal("API key not found in path")
		case r.URL.Query().Get("chat_id") != cfg.ChatID:
			t.Fatal("Invalid chat_id")
		case !strings.Contains(r.URL.Query().Get("text"), addText):
			t.Fatal("Add text not found")
		case !strings.Contains(r.URL.Query().Get("text"), description):
			t.Fatal("Description not found")
		case !strings.Contains(r.URL.Query().Get("text"), username):
			t.Fatal("Username not found")
		}
	}))
	defer ts.Close()

	cfg.URL = ts.URL

	tm := impl.NewTelegramNotifier(cfg)
	err := tm.Send(
		context.Background(),
		map[string]string{
			glmt.TmpVarDescription: description,
		},
		addText,
		[]*team.Member{{
			Names: map[string]string{
				"telegram_member_id": username,
			},
		}},
	)

	select {
	case <-tsDone:
	case <-time.After(time.Second):
		t.Fatal("Server response timeout")
	}

	if err != nil {
		t.Fatal(err)
	}
}
