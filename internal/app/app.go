package app

import (
	"context"
	"fmt"
	"os"

	"go-ldap-slack-syncer/internal/actions/users"
	"go-ldap-slack-syncer/internal/config"
	"go-ldap-slack-syncer/internal/infra/database"
	ldapsync "go-ldap-slack-syncer/internal/infra/ldap"
	slacksync "go-ldap-slack-syncer/internal/infra/slack"
	"go-ldap-slack-syncer/internal/logging"
	"go-ldap-slack-syncer/internal/services"
)

func RunApp() {
	var (
		debug bool
	)

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	log, err := logging.Configure(debug)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		os.Exit(2)
	}

	conf, flagsConf, err := config.Parse(os.Args)
	if err != nil {
		log.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Time.MaximumWorking)
	defer cancel()

	db, err := database.InitDatabase(ctx, conf.Storage.MySQL)
	if err != nil {
		log.Error(err)
		return
	}

	ctx = database.ContextWithDatabase(ctx, db)

	ldapClient, err := ldapsync.GetClient(ctx, conf.Ldap.Host, conf.Ldap.Account, conf.Ldap.Settings, conf.Ldap.Security)
	if err != nil {
		log.Error(err)
		return
	}

	defer ldapClient.Connection.Close()

	services.InitUsers(ldapClient, slacksync.GetClient(conf.Slack))

	usersActionService := users.GetActionService(services.GetUsersService())

	if flagsConf.Revert {
		if flagsConf.Disable {
			err = usersActionService.DisableUsersFromDB(ctx, flagsConf.RevertDate, flagsConf.UserMail, flagsConf.Apply)
			if err != nil {
				log.Error(err)
				return
			}
		}

		if flagsConf.Enable {
			err = usersActionService.EnableUsersFromDB(ctx, flagsConf.RevertDate, flagsConf.UserMail, flagsConf.Apply)
			if err != nil {
				log.Error(err)
				return
			}
		}

		return
	}

	usersDeletionQueue, err := usersActionService.FindDeletedInSlack(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	err = usersActionService.DisableNewUsers(ctx, usersDeletionQueue, flagsConf.Apply)
	if err != nil {
		log.Error(err)
		return
	}
}
