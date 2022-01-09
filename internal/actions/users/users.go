package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/config"
	"go-ldap-slack-syncer/internal/infra/database/dbservice"
	dbtypes "go-ldap-slack-syncer/internal/infra/database/dbservice/types"
	"go-ldap-slack-syncer/internal/infra/slack/actions"
	"go-ldap-slack-syncer/internal/logging"
	"go-ldap-slack-syncer/internal/services"
	"go-ldap-slack-syncer/internal/services/types"
	"go-ldap-slack-syncer/internal/sync"
)

type ActionService struct {
	usersService types.Service
}

func (s ActionService) FindDeletedInSlack(ctx context.Context) ([]slack.User, error) {
	var (
		log = logging.GetLogger()
	)

	usersDeletionQueue, err := sync.CompareUsers(ctx, s.usersService)
	if err != nil {
		return nil, err
	}

	log.Infof("total number of users found to be deleted in slack workspace: %d", len(usersDeletionQueue))

	return usersDeletionQueue, nil
}

func (s ActionService) DisableNewUsers(ctx context.Context, users []slack.User, apply bool) error {
	var (
		client = s.usersService.GetSlackClient(ctx)

		flagsConf = config.GetFlagsConfig()

		disabledUsersCount int
		messages           []string
		usersForDisable    []slack.User
	)

	for _, user := range users {
		userInfoFromDB, err := dbservice.GetUser(ctx, user.ID, user.Profile.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			messages = append(messages, fmt.Sprintf("User with mail '%s' and slackID '%s' has already been processed by the utility on %s. Skipping...\n", user.Profile.Email, user.ID, userInfoFromDB.Action))

			continue
		}

		if apply {
			usersForDisable = append(usersForDisable, user)
			disabledUsersCount++
		} else {
			messages = append(messages, fmt.Sprintf("User with mail '%s' and slackID '%s' can be disabled in slack space %s. Last updated: %s\n", user.Profile.Email, user.ID, client.TeamName, user.Updated.Time().Format("02-Jan-2006")))

			disabledUsersCount++
		}

		if flagsConf.Count != 0 {
			if disabledUsersCount == int(flagsConf.Count) {
				break
			}
		}
	}

	if apply {
		if len(messages) > 0 {
			_, err := actions.SendMessages(ctx, client, messages, true)
			if err != nil {
				return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
			}
		}

		err := actions.DisableNewUsers(ctx, client, usersForDisable)
		if err != nil {
			return fmt.Errorf("error disabled new users: %w", err)
		}
	} else {
		messages = append(messages, fmt.Sprintf("Total number of users in the disabled queue in slack: %d", disabledUsersCount))

		_, err := actions.SendMessages(ctx, client, messages, true)
		if err != nil {
			return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
		}
	}

	return nil
}

func (s ActionService) DisableUsersFromDB(ctx context.Context, date time.Time, userMail string, apply bool) error {
	var (
		client = s.usersService.GetSlackClient(ctx)

		flagsConf = config.GetFlagsConfig()

		disabledUsersCount int
		messages           []string
		usersForDisable    []slack.User

		found bool
	)

	dbUsers, err := dbservice.GetUsersForDate(ctx, date, dbtypes.UserEnable)
	if err != nil {
		return err
	}

	for _, dbUser := range dbUsers {
		if userMail != "" {
			if found {
				break
			}

			if dbUser.Mail != userMail {
				continue
			}

			found = true
		}

		if apply {
			usersForDisable = append(usersForDisable, slack.User{
				ID: dbUser.SlackID,
				Profile: slack.UserProfile{
					Email: dbUser.Mail,
				},
			})

			disabledUsersCount++
		} else {
			messages = append(messages, fmt.Sprintf("User with mail '%s' and slackID '%s' can be disabled in slack space %s\n", dbUser.Mail, dbUser.SlackID, client.TeamName))

			disabledUsersCount++
		}

		if flagsConf.Count != 0 {
			if disabledUsersCount == int(flagsConf.Count) {
				break
			}
		}
	}

	if userMail != "" {
		if len(usersForDisable) == 0 && len(messages) == 0 {
			return fmt.Errorf("the specified user with mail %s was not found in the database", userMail)
		}
	}

	if apply {
		if len(messages) > 0 {
			_, err = actions.SendMessages(ctx, client, messages, true)
			if err != nil {
				return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
			}
		}

		err = actions.DisableUsers(ctx, client, usersForDisable)
		if err != nil {
			return fmt.Errorf("error disabled new users: %w", err)
		}
	} else {
		messages = append(messages, fmt.Sprintf("Total number of users in the disabled queue in slack: %d", disabledUsersCount))

		_, err = actions.SendMessages(ctx, client, messages, true)
		if err != nil {
			return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
		}
	}

	return nil
}

func (s ActionService) EnableUsersFromDB(ctx context.Context, date time.Time, userMail string, apply bool) error {
	var (
		client = s.usersService.GetSlackClient(ctx)

		flagsConf = config.GetFlagsConfig()

		enabledUsersCount int
		messages          []string
		usersForEnable    []slack.User

		found bool
	)

	dbUsers, err := dbservice.GetUsersForDate(ctx, date, dbtypes.UserDisable)
	if err != nil {
		return err
	}

	for _, dbUser := range dbUsers {
		if userMail != "" {
			if found {
				break
			}

			if dbUser.Mail != userMail {
				continue
			}

			found = true
		}

		if apply {
			usersForEnable = append(usersForEnable, slack.User{
				ID: dbUser.SlackID,
				Profile: slack.UserProfile{
					Email: dbUser.Mail,
				},
			})

			enabledUsersCount++
		} else {
			messages = append(messages, fmt.Sprintf("User with mail '%s' and slackID '%s' can be enabled in slack space %s\n", dbUser.Mail, dbUser.SlackID, client.TeamName))

			enabledUsersCount++
		}

		if flagsConf.Count != 0 {
			if enabledUsersCount == int(flagsConf.Count) {
				break
			}
		}

	}

	if userMail != "" {
		if len(usersForEnable) == 0 && len(messages) == 0 {
			return fmt.Errorf("the specified user with mail %s was not found in the database", userMail)
		}
	}

	if apply {
		if len(messages) > 0 {
			_, err = actions.SendMessages(ctx, client, messages, true)
			if err != nil {
				return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
			}
		}

		err = actions.EnableUsers(ctx, client, usersForEnable)
		if err != nil {
			return fmt.Errorf("error enabled new users: %w", err)
		}
	} else {
		messages = append(messages, fmt.Sprintf("Total number of users in the enabled queue in slack: %d", enabledUsersCount))

		_, err = actions.SendMessages(ctx, client, messages, true)
		if err != nil {
			return fmt.Errorf("error send messages to channel '%s': %s", client.NotificationChannelID, err)
		}
	}

	return nil
}

func GetActionService(users services.Users) ActionService {
	return ActionService{usersService: users}
}
