package sync

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/config"
	"go-ldap-slack-syncer/internal/services/types"
)

var (
	ErrInvalidDataType = errors.New("invalid data type")
)

func CompareUsers(ctx context.Context, service types.Service) ([]slack.User, error) {
	var (
		flagsConf = config.GetFlagsConfig()
	)

	ldapData, err := service.GetLdapData(ctx)
	if err != nil {
		return nil, err
	}

	slackData, err := service.GetSlackData(ctx)
	if err != nil {
		return nil, err
	}

	ldapUsers, ok := ldapData.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("%w: expected map[string]string", ErrInvalidDataType)
	}

	slackUsers, ok := slackData.([]slack.User)
	if !ok {
		return nil, fmt.Errorf("%w: expected []slack.User", ErrInvalidDataType)
	}

	usersDeletionQueue := compare(ctx, flagsConf, ldapUsers, slackUsers)

	return usersDeletionQueue, nil
}

// test
func compare(ctx context.Context, flagsConf config.FlagsConfig, ldapUsers map[string]string, slackUsers []slack.User) []slack.User {
	var (
		usersDeletionQueue []slack.User
	)

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

	return usersDeletionQueue
}
