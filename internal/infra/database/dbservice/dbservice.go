package dbservice

import (
	"context"
	"fmt"
	"time"

	"go-ldap-slack-syncer/internal/infra/database"
	dbtypes "go-ldap-slack-syncer/internal/infra/database/dbservice/types"
	"go-ldap-slack-syncer/internal/logging"
)

func AddUser(ctx context.Context, user dbtypes.User) error {
	var (
		log = logging.GetLogger()
	)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return err
	}

	log.Debugf("inserting user with mail %s and action '%s' to database", user.Mail, user.Action)

	_, err = db.GetDBConn().ExecContext(ctx, "insert into users (slack_id, mail, action, date_of_action) values (?, ?, ?, ?)", user.SlackID, user.Mail, user.Action, user.Date)
	if err != nil {
		return fmt.Errorf("error inserting user with mail %s and action '%s' into the database: %w", user.Mail, user.Action, err)
	}

	return nil
}

func GetUser(ctx context.Context, slackID, mail string) (dbtypes.User, error) {
	var (
		log = logging.GetLogger()

		result dbtypes.User
	)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return dbtypes.User{}, err
	}

	log.Debugf("getting user with mail %s from database", mail)

	row := db.GetDBConn().QueryRowContext(ctx, "select slack_id, mail, action, date_of_action from users where slack_id = ? and mail = ?", slackID, mail)

	err = row.Scan(&result.SlackID, &result.Mail, &result.Action, &result.Date)
	if err != nil {
		return dbtypes.User{}, fmt.Errorf("error getting a user with mail %s from the database: %w", mail, err)
	}

	return result, nil
}

func GetUsersForDate(ctx context.Context, date time.Time, action dbtypes.UserAction) ([]dbtypes.User, error) {
	var (
		log = logging.GetLogger()

		result []dbtypes.User
	)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("getting users for date %s from database", date.String())

	rows, err := db.GetDBConn().QueryContext(ctx, "select slack_id, mail, action, date_of_action from users where date_of_action = ? and action = ?", date, action)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Errorf("the attempt to close rows was unsuccessful: %s", err.Error())
		}
	}()

	for rows.Next() {
		var user dbtypes.User

		err = rows.Scan(&user.SlackID, &user.Mail, &user.Action, &user.Date)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func GetUsersForAction(ctx context.Context, action dbtypes.UserAction) ([]dbtypes.User, error) {
	var (
		log = logging.GetLogger()

		result []dbtypes.User
	)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("getting users for action %s from database", action)

	rows, err := db.GetDBConn().QueryContext(ctx, "select slack_id, mail, action, date_of_action from users where action = ?", action)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Errorf("the attempt to close rows was unsuccessful: %s", err.Error())
		}
	}()

	for rows.Next() {
		var user dbtypes.User

		err = rows.Scan(&user.SlackID, &user.Mail, &user.Action, &user.Date)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func UpdateUserAction(ctx context.Context, user dbtypes.User) error {
	var (
		log = logging.GetLogger()
	)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return err
	}

	log.Debugf("update user with mail %s to action %s in database", user.Mail, user.Action)

	_, err = db.GetDBConn().ExecContext(ctx, "update users set action = ?, date_of_action = ? where slack_id = ? and mail = ?", user.Action, user.Date, user.SlackID, user.Mail)
	if err != nil {
		return fmt.Errorf("error update user with mail %s to action '%s' in database: %w", user.Mail, user.Action, err)
	}

	return nil
}
