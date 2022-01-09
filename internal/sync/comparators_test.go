package sync

import (
	"context"
	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/config"
	"testing"
)

func TestCompare(t *testing.T) {
	var (
		ldapUsers  = map[string]string{"a@mail.ru": "", "b@mail.ru": ""}
		slackUsers = []slack.User{{
			Profile: slack.UserProfile{Email: "a@mail.ru"},
		}, {
			Profile: slack.UserProfile{Email: "B@mail.ru"},
		}}
	)

	got := compare(context.Background(), config.FlagsConfig{}, ldapUsers, slackUsers)

	if got != nil {
		t.Errorf("got %d, wanted nil", len(got))
	}
}
