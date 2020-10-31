package impl

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"
)

func TestHTTPMentioner(t *testing.T) {
	const (
		currentMemberUID = "00000000-0000-0000-0000-000000000000"
		currentProject   = "hummerd/glmt"
	)

	m := mentioner.Member{IsActive: true}
	expMembers := []mentioner.Member{m, m}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mentionsConfig := mentionsConfig{
			Members: expMembers,
		}

		err := json.NewEncoder(w).Encode(&mentionsConfig)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	tsURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var httpMentioner mentioner.Mentioner = NewHTTPMentioner(
		tsURL,
		currentMemberUID,
		len(expMembers),
	)

	var gotMembers []mentioner.Member
	gotMembers, err = httpMentioner.Mentions(currentProject)
	switch {
	case err != nil:
		t.Fatal(err)
	case len(gotMembers) != len(expMembers):
		t.Fatalf("len(members) exp: %d, got: %d", len(expMembers), len(gotMembers))
	}
}
