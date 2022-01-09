package types

import (
	"context"

	ldapinfratypes "go-ldap-slack-syncer/internal/infra/ldap/types"
	slackinfratypes "go-ldap-slack-syncer/internal/infra/slack/types"
)

type Service interface {
	GetLdapClient(ctx context.Context) *ldapinfratypes.Client
	GetSlackClient(ctx context.Context) *slackinfratypes.Client
	GetLdapData(ctx context.Context) (data interface{}, err error)
	GetSlackData(ctx context.Context) (data interface{}, err error)
}
