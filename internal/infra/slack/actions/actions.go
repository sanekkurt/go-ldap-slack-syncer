package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"go-ldap-slack-syncer/internal/infra/database/dbservice"
	dbtypes "go-ldap-slack-syncer/internal/infra/database/dbservice/types"
	"go-ldap-slack-syncer/internal/infra/slack/constants"
	"go-ldap-slack-syncer/internal/infra/slack/types"
	slackutils "go-ldap-slack-syncer/internal/infra/slack/utils"
	"go-ldap-slack-syncer/internal/logging"
)

func GetUsers(ctx context.Context, connect *slack.Client) ([]slack.User, error) {
	var (
		log = logging.GetLogger()
	)

	log.Info("getting slack users...")

	users, err := connect.GetUsersContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("the total number of users received from the slack api: %d", len(users))

	return users, nil
}

func SendMessages(ctx context.Context, client *types.Client, messages []string, isJoin bool) (timestamp string, err error) {
	var (
		log = logging.GetLogger()
	)

	if isJoin {
		messages = slackutils.SplitStr(messages, 4000)
	}

	for _, message := range messages {
		log.Debugf("send message '%s' to channel '%s'", message, client.NotificationChannelID)

		_, timestamp, err = client.BotUserClient.PostMessageContext(ctx, client.NotificationChannelID, slack.MsgOptionText(message, false))
		if err != nil {
			return "", err
		}

		time.Sleep(constants.TimeoutBetweenPostingMessages)
	}

	return
}

func sendMessage(ctx context.Context, client *types.Client, message string) (timestamp string, err error) {
	var (
		log = logging.GetLogger()
	)

	if len(message) > constants.LimitCharacters {
		log.Infof("The message size exceeds 4000 characters so it will be divided into several messages")
	}

	log.Debugf("send message '%s' to channel '%s'", message, client.NotificationChannelID)

	_, timestamp, err = client.BotUserClient.PostMessageContext(ctx, client.NotificationChannelID, slack.MsgOptionText(message, false))
	if err != nil {
		return "", err
	}

	return
}

func UpdateUserRealName(ctx context.Context, connect *slack.Client, userID, firstName, lastName string) error {
	var (
		log = logging.GetLogger()
	)

	log.Infof("update real name for slack user '%s' to %s %s", userID, firstName, lastName)

	err := connect.SetUserRealNameContextWithUser(ctx, userID, fmt.Sprintf("%s %s", firstName, lastName))
	if err != nil {
		return err
	}

	return nil
}

func DisableUsers(ctx context.Context, client *types.Client, users []slack.User) error {
	var (
		log = logging.GetLogger()

		disabledUsersCount int
	)

	for _, user := range users {
		log.Infof("disable slack user '%s' in team '%s'", user.ID, client.TeamName)

		// TODO ЛОГИКА ВЫКЛЮЧЕНИЯ ЮЗЕРА
		//err := client.UserClient.DisableUserContext(ctx, client.TeamName, user.ID)
		//if err != nil {
		//	_, _ = sendMessage(ctx, client, fmt.Sprintf("Error disable user with mail '%s' and slackID '%s': %s. Stopping the disabled process...",user.ID, user.Profile.Email, err.Error()))
		//	_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))
		//
		//	return err
		//}

		disabledUsersCount++

		err := dbservice.UpdateUserAction(ctx, dbtypes.User{
			SlackID: user.ID,
			Mail:    user.Profile.Email,
			Action:  dbtypes.UserDisable,
			Date:    time.Now()})
		if err != nil {
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Error adding a disabled user to the database: %s", err.Error()))
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))

			return err
		}

		_, err = sendMessage(ctx, client, fmt.Sprintf("User with mail '%s' and slackID '%s' disable in slack space", user.Profile.Email, user.ID))
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Minute / constants.Tier2)
	}

	_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))

	return nil
}

func DisableNewUsers(ctx context.Context, client *types.Client, users []slack.User) error {
	var (
		log = logging.GetLogger()

		disabledUsersCount int
	)

	for _, user := range users {
		log.Infof("disable slack user '%s' in team '%s'", user.ID, client.TeamName)

		//TODO ЛОГИКА ВЫКЛЮЧЕНИЯ ЮЗЕРА
		err := client.UserClient.DisableUserContext(ctx, client.TeamName, user.ID)

		if err != nil {
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Error disable user with mail '%s' and slackID '%s': %s. Stopping the disabled process...", user.ID, user.Profile.Email, err.Error()))
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))

			return err
		}

		disabledUsersCount++

		err = dbservice.AddUser(ctx, dbtypes.User{
			SlackID: user.ID,
			Mail:    user.Profile.Email,
			Action:  dbtypes.UserDisable,
			Date:    time.Now()})
		if err != nil {
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Error adding a disabled user to the database: %s", err.Error()))
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))

			return err
		}

		_, err = sendMessage(ctx, client, fmt.Sprintf("User with mail '%s' and slackID '%s' disable in slack space", user.Profile.Email, user.ID))
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Minute / constants.Tier2)
	}

	_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", disabledUsersCount))

	return nil
}

func EnableUsers(ctx context.Context, client *types.Client, users []slack.User) error {
	var (
		log = logging.GetLogger()

		enabledUsersCount int
	)

	for _, user := range users {
		log.Infof("enable slack user '%s' in team '%s'", user.ID, client.TeamName)

		// TODO ЛОГИКА ВКЛЮЧЕНИЯ ЮЗЕРА
		err := client.UserClient.SetRegularContext(ctx, client.TeamName, user.ID)
		if err != nil {
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Error enable user with mail '%s' and slackID '%s': %s. Stopping the enabled process...", user.ID, user.Profile.Email, err.Error()))
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of disabled users in the slack: %d", enabledUsersCount))

			return err
		}

		enabledUsersCount++

		err = dbservice.UpdateUserAction(ctx, dbtypes.User{
			SlackID: user.ID,
			Mail:    user.Profile.Email,
			Action:  dbtypes.UserEnable,
			Date:    time.Now()})
		if err != nil {
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Error adding a enabled user to the database: %s", err.Error()))
			_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of enabled users in the slack: %d", enabledUsersCount))

			return err
		}

		_, err = sendMessage(ctx, client, fmt.Sprintf("User with mail '%s' and slackID '%s' enable in slack space", user.Profile.Email, user.ID))
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Minute / constants.Tier2)
	}

	_, _ = sendMessage(ctx, client, fmt.Sprintf("Total number of enabled users in the slack: %d", enabledUsersCount))

	return nil
}
