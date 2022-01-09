package types

import (
	"time"

	"github.com/go-ldap/ldap/v3"
)

type LDAPInputData struct {
	Host
	Account
	Settings
}

type Host struct {
	Address string `yaml:"address"`
	Port    uint   `yaml:"port"`
}

type Account struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Settings struct {
	TimeLimitConnection time.Duration `yaml:"timeLimitConnection"`
	PagingSize          uint32        `yaml:"pagingSize"`
	SizeLimit           uint          `yaml:"sizeLimit"`
	BaseDNSearchSuffix  string        `yaml:"baseDnSearchSuffix"`
}

type Security struct {
	MinNumberUsers uint `yaml:"minNumberUsers"`
}

type Client struct {
	Connection *ldap.Conn
	Account
	Settings
	Security
}
