package sync

import (
	"context"
	"strings"

	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/config"
	"go-ldap-slack-syncer/internal/services/types"
)

func CompareUsers(ctx context.Context, service types.Service) ([]slack.User, error) {
	var (
		usersDeletionQueue []slack.User
		flagsConf          = config.GetFlagsConfig()
	)

	ldapData, err := service.GetLdapData(ctx)
	if err != nil {
		return nil, err
	}

	slackData, err := service.GetSlackData(ctx)
	if err != nil {
		return nil, err
	}

	ldapUsers := ldapData.(map[string]string)
	slackUsers := slackData.([]slack.User)
	for _, slackUser := range slackUsers {
		if slackUser.Deleted || slackUser.IsBot {
			continue
		}

		if _, ok := ldapUsers[strings.ToLower(slackUser.Profile.Email)]; !ok {
			if flagsConf.LastUpdate != nil {
				if flagsConf.BeforeLastUpdate {
					if !slackUser.Updated.Time().Before(*flagsConf.LastUpdate) {
						continue
					}
				}

				if flagsConf.AfterLastUpdate {
					if !slackUser.Updated.Time().After(*flagsConf.LastUpdate) {
						continue
					}
				}
			}
			usersDeletionQueue = append(usersDeletionQueue, slackUser)
		}
	}

	return usersDeletionQueue, nil
}
