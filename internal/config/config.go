package config

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
)

var (
	opts struct {
		ConfigPath       string `long:"config" short:"c" env:"CONFIG_PATH" description:"Path to config.yaml file" required:"true"`
		Apply            bool   `long:"apply" description:"enable apply mod for service"`
		Revert           bool   `long:"revert" description:"enable revert mod for service"`
		RevertDate       string `long:"date" short:"d" description:"date for revert mod for service"`
		Enable           bool   `long:"enable" description:"enable users for revert mode"`
		Disable          bool   `long:"disable" description:"disable users for revert mode"`
		UserMail         string `long:"usermail" short:"u" description:"revert action only for the specified user"`
		Count            uint   `long:"count" description:"the maximum number of users that can be changed"`
		LastUpdate       string `long:"lastupdate" description:"date of the last update of the user in the slack"`
		BeforeLastUpdate bool   `long:"before" description:"synchronize users before the last update date"`
		AfterLastUpdate  bool   `long:"after" description:"synchronize users after the last update date"`
		//SyncUserName bool   `long:"syncusername" description:"synchronize ldap user names with slack"`
	}
	ErrHelpShown = errors.New("help message shown")
)

func Parse(args []string) (*AppConfig, *FlagsConfig, error) {
	_, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).ParseArgs(args[1:])
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				return nil, nil, ErrHelpShown
			}

			return nil, nil, fmt.Errorf("cannot parse arguments: %w", flagsErr)
		}

		return nil, nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	cfg, err := ParseConfig(opts.ConfigPath)
	if err != nil {
		return nil, nil, err
	}

	flagsCfg, err := ParseFlagsConfig(opts.Apply, opts.Revert, opts.Enable, opts.Disable, opts.RevertDate, opts.UserMail, opts.Count, opts.LastUpdate, opts.BeforeLastUpdate, opts.AfterLastUpdate)
	if err != nil {
		return nil, nil, err
	}

	return cfg, flagsCfg, nil
}
