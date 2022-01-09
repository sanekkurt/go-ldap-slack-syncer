package ldapsync

import (
	"fmt"
	"strconv"

	"go-ldap-slack-syncer/internal/infra/ldap/types"
)

func createURL(host types.Host) string {
	return fmt.Sprintf("ldap://" + host.Address + ":" + strconv.Itoa(int(host.Port)))
}

//func getBaseDNFromString(input string) string {
//	var (
//		res   []string
//		parts = strings.Split(input, ",")
//	)
//
//	for _, part := range parts {
//		sign := strings.Split(part, "=")
//		if len(sign) == 2 {
//			if sign[0] == "dc" || sign[0] == "DC" {
//				res = append(res, part)
//			}
//		}
//	}
//
//	return strings.Join(res, ",")
//}

func generateUserAttributesFilter(attributes map[string]string) string {
	var filter string

	for key, value := range attributes {
		filter += "(" + key + "=" + value + ")"
	}

	if filter != "" {
		filter = "(&(objectClass=user)" + filter + ")"
	}

	return filter
}

func attributesMapToSlice(attributes map[string]string) []string {
	var (
		attributesCount = len(attributes)
		result          = make([]string, attributesCount)
		index           = 0
	)

	for attributeName := range attributes {
		result[index] = attributeName

		index++
	}

	return result
}
