package configstructs

import (
	ldapinfratypes "go-ldap-slack-syncer/internal/infra/ldap/types"
)

type Ldap struct {
	ldapinfratypes.Host     `yaml:"host"`
	ldapinfratypes.Account  `yaml:"account"`
	ldapinfratypes.Settings `yaml:"settings"`
	ldapinfratypes.Security `yaml:"security"`
}
