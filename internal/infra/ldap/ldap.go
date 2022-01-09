package ldapsync

import (
	"context"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	ldapinfratypes "go-ldap-slack-syncer/internal/infra/ldap/types"
	"go-ldap-slack-syncer/internal/logging"
)

func GetClient(ctx context.Context, host ldapinfratypes.Host, account ldapinfratypes.Account, settings ldapinfratypes.Settings, security ldapinfratypes.Security) (*ldapinfratypes.Client, error) {
	l, err := ldap.DialURL(createURL(host))
	if err != nil {
		return nil, err
	}
	// defer l.Close()
	err = l.Bind(account.Username, account.Password)
	if err != nil {
		return nil, err
	}

	return &ldapinfratypes.Client{Connection: l, Account: account, Settings: settings, Security: security}, nil
}

func GetActiveUsersWithMail(ctx context.Context, client *ldapinfratypes.Client) ([]*ldap.Entry, error) {
	users, err := findUsersWithTheAttributes(ctx, client, map[string]string{"mail": "*", "name": "*"})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func findUsersWithTheAttributes(ctx context.Context, client *ldapinfratypes.Client, attributes map[string]string) ([]*ldap.Entry, error) {
	var (
		log = logging.GetLogger()

		filter        = generateUserAttributesFilter(attributes)
		searchRequest = ldap.NewSearchRequest(client.BaseDNSearchSuffix, ldap.ScopeWholeSubtree, 0, int(client.SizeLimit), int(client.TimeLimitConnection.Seconds()), false, filter, attributesMapToSlice(attributes), nil)
	)

	log.Info("getting ldap users...")

	sr, err := client.Connection.SearchWithPaging(searchRequest, client.PagingSize)
	if err != nil {
		return nil, err
	}

	usersCount := len(sr.Entries)

	log.Debugf("the total number of users received from the ldap api: %d", usersCount)

	if uint(usersCount) < client.Security.MinNumberUsers {
		return nil, fmt.Errorf("the number of users received from the LDAP server is less than the safe minimum number of users specified: %d", client.Security.MinNumberUsers)
	}

	return sr.Entries, nil
}
