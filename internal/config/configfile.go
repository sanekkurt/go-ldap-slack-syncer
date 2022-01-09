package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-ldap-slack-syncer/internal/config/configstructs"
	slackinfratypes "go-ldap-slack-syncer/internal/infra/slack/types"
	"go-ldap-slack-syncer/internal/logging"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Ldap    configstructs.Ldap             `yaml:"ldap"`
	Slack   slackinfratypes.SlackInputData `yaml:"slack"`
	Storage configstructs.Storage          `yaml:"storage"`
	Time    configstructs.Time             `yaml:"time"`
}

type FlagsConfig struct {
	Apply            bool
	Revert           bool
	RevertDate       time.Time
	Enable           bool
	Disable          bool
	UserMail         string
	Count            uint
	LastUpdate       *time.Time
	BeforeLastUpdate bool
	AfterLastUpdate  bool
	//SyncUserName bool
}

var globalFlagsConfig *FlagsConfig

const (
	defaultLdapPort                uint   = 389
	defaultLdapPagingSize          uint32 = 50
	defaultLdapSizeLimit           uint   = 100000
	defaultLdapTimeLimitConnection        = 1 * time.Minute
	defaultMySQLPort                      = 3306
	defaultMaxWorkingApp                  = 1 * time.Hour
)

func ParseConfig(cfgPath string) (*AppConfig, error) {
	log := logging.GetLogger()

	configFullPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path to %s: %w", configFullPath, err)
	}

	f, err := os.Open(configFullPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("cannot open config '%s' for reading: %w", configFullPath, err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Warnf("cannot close config '%s': %s", cfgPath, err)
		}
	}()

	dec := yaml.NewDecoder(f)
	dec.SetStrict(true)

	cfg := &AppConfig{}

	err = dec.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config from '%s': %w", cfgPath, err)
	}

	if cfg.Slack.NotificationChannelID == "" {
		return nil, fmt.Errorf("error parse config: slack.notificationChannelID must be specified")
	}

	if cfg.Slack.BotUserOAuthToken == "" {
		return nil, fmt.Errorf("error parse config: slack.botUserOAuthToken must be specified")
	}

	if cfg.Slack.UserOAuthToken == "" {
		return nil, fmt.Errorf("error parse config: slack.userOAuthToken must be specified")
	}

	if cfg.Slack.WorkspaceName == "" {
		return nil, fmt.Errorf("error parse config: slack.workspaceName must be specified")
	}

	if cfg.Ldap.Address == "" {
		return nil, fmt.Errorf("error parse config: ldap.host.address must be specified")
	}

	if cfg.Ldap.Port == 0 {
		log.Infof("the ldap.host.port was not specified in the config file. Installed by default: %d", defaultLdapPort)

		cfg.Ldap.Port = defaultLdapPort
	}

	if cfg.Ldap.Account.Username == "" {
		return nil, fmt.Errorf("error parse config: ldap.account.username must be specified")
	}

	if cfg.Ldap.Account.Password == "" {
		return nil, fmt.Errorf("error parse config: ldap.account.password must be specified")
	}

	if cfg.Ldap.Settings.PagingSize == 0 {
		log.Infof("the ldap.settings.pagingSize was not specified in the config file. Installed by default: %d", defaultLdapPagingSize)

		cfg.Ldap.Settings.PagingSize = defaultLdapPagingSize
	}

	if cfg.Ldap.Settings.SizeLimit == 0 {
		log.Infof("the ldap.settings.sizeLimit was not specified in the config file. Installed by default: %d", defaultLdapSizeLimit)

		cfg.Ldap.Settings.SizeLimit = defaultLdapSizeLimit
	}

	if cfg.Ldap.Settings.TimeLimitConnection == 0 {
		log.Infof("the ldap.settings.timeLimitConnection was not specified in the config file. Installed by default: %s", defaultLdapTimeLimitConnection.String())

		cfg.Ldap.Settings.TimeLimitConnection = defaultLdapTimeLimitConnection
	}

	if cfg.Ldap.BaseDNSearchSuffix == "" {
		return nil, fmt.Errorf("error parse config: ldap.settings.baseDnSearchSuffix must be specified")
	}

	if cfg.Storage.MySQL.Address == "" {
		return nil, fmt.Errorf("error parse config: storage.mysql.address must be specified")
	}

	if cfg.Storage.MySQL.Database == "" {
		return nil, fmt.Errorf("error parse config: storage.mysql.database must be specified")
	}

	if cfg.Storage.MySQL.Username == "" {
		return nil, fmt.Errorf("error parse config: storage.mysql.username must be specified")
	}

	if cfg.Storage.MySQL.Password == "" {
		return nil, fmt.Errorf("error parse config: storage.mysql.password must be specified")
	}

	if cfg.Storage.MySQL.Port == 0 {
		log.Infof("the storage.mysql.port was not specified in the config file. Installed by default: %d", defaultMySQLPort)

		cfg.Storage.MySQL.Port = defaultMySQLPort
	}

	if cfg.Time.MaximumWorking == 0 {
		log.Infof("the time.maximumWorking was not specified in the config file. Installed by default: %s", defaultMaxWorkingApp.String())

		cfg.Time.MaximumWorking = defaultMaxWorkingApp
	}

	return cfg, nil
}

func ParseFlagsConfig(apply, revert, enable, disable bool, revertDate, mailUser string, count uint, lastUpdate string, beforeLastUpdate, afterLastUpdate bool) (*FlagsConfig, error) {
	var cfg FlagsConfig

	if revert {
		if revertDate == "" {
			return nil, fmt.Errorf("incorrect configuration of flags. In order to enable reverse mode, you need to pass the -d or --date flag in the YYYY-MM-DD format to indicate the day when the user status is restored")
		}

		if enable && disable {
			return nil, fmt.Errorf("incorrect configuration of flags. You have to choose only one action for revert. enable or disable")
		}

		if !enable && !disable {
			return nil, fmt.Errorf("incorrect configuration of flags. You have to choose one action for revert. enable or disable")
		}

		if mailUser != "" {
			cfg.UserMail = mailUser
		}

		if revertDate != "" {
			date, err := time.Parse("2006-01-02", revertDate)
			if err != nil {
				return nil, err
			}

			cfg.RevertDate = date
		}

		cfg.Disable = disable
		cfg.Enable = enable
	}

	cfg.Apply = apply
	cfg.Revert = revert
	cfg.Count = count

	if lastUpdate != "" {
		if beforeLastUpdate && afterLastUpdate {
			return nil, fmt.Errorf("incorrect configuration of flags. You have to choose only one action for last update. before or after")
		}

		if !beforeLastUpdate && !afterLastUpdate {
			return nil, fmt.Errorf("incorrect configuration of flags. You have to choose one action for last update. before or after")
		}

		date, err := time.Parse("2006-01-02", lastUpdate)
		if err != nil {
			return nil, err
		}

		cfg.LastUpdate = &date
		cfg.BeforeLastUpdate = beforeLastUpdate
		cfg.AfterLastUpdate = afterLastUpdate
	}

	globalFlagsConfig = &cfg

	return globalFlagsConfig, nil
}

func GetFlagsConfig() FlagsConfig {
	if globalFlagsConfig == nil {
		return FlagsConfig{}
	}

	return *globalFlagsConfig
}
