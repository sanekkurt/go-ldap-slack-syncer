package types

import (
	"github.com/slack-go/slack"
)

type BotUserToken string
type UserToken string

type SlackInputData struct {
	UserOAuthToken        UserToken    `yaml:"userOAuthToken"`
	BotUserOAuthToken     BotUserToken `yaml:"botUserOAuthToken"`
	NotificationChannelID string       `yaml:"notificationChannelID"`
	WorkspaceName         string       `yaml:"workspaceName"`
}

type Client struct {
	UserClient            *slack.Client
	BotUserClient         *slack.Client
	NotificationChannelID string
	TeamName              string
}
