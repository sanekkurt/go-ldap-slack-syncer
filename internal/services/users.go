package services

import (
	"context"
	"strings"

	ldapsync "go-ldap-slack-syncer/internal/infra/ldap"
	ldapinfratypes "go-ldap-slack-syncer/internal/infra/ldap/types"
	"go-ldap-slack-syncer/internal/infra/slack/actions"
	slackinfratypes "go-ldap-slack-syncer/internal/infra/slack/types"
)

type Users struct {
	ldapClient  *ldapinfratypes.Client
	slackClient *slackinfratypes.Client
}

var u Users

func InitUsers(l *ldapinfratypes.Client, s *slackinfratypes.Client) {
	u.ldapClient = l
	u.slackClient = s
}

func GetUsersService() Users {
	return u
}

func (u Users) GetLdapClient(ctx context.Context) *ldapinfratypes.Client {
	return u.ldapClient
}

func (u Users) GetSlackClient(ctx context.Context) *slackinfratypes.Client {
	return u.slackClient
}

func (u Users) GetLdapData(ctx context.Context) (data interface{}, err error) {
	users, err := ldapsync.GetActiveUsersWithMail(ctx, u.ldapClient)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for _, user := range users {
		mail := user.GetAttributeValue("mail")
		name := user.GetAttributeValue("name")
		result[strings.ToLower(mail)] = name
	}

	return result, nil
}

func (u Users) GetSlackData(ctx context.Context) (data interface{}, err error) {
	result, err := actions.GetUsers(ctx, u.slackClient.UserClient)
	if err != nil {
		return nil, err
	}

	return result, nil
}
