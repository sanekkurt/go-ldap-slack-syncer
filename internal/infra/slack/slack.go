package slacksync

import (
	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/infra/slack/types"
)

func GetClient(data types.SlackInputData) *types.Client {
	return &types.Client{
		UserClient:            slack.New(string(data.UserOAuthToken)),
		BotUserClient:         slack.New(string(data.BotUserOAuthToken)),
		NotificationChannelID: data.NotificationChannelID,
		TeamName:              data.WorkspaceName,
	}
}
