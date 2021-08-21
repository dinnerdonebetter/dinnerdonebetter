package bleve

import (
	"fmt"
	"regexp"
)

var (
	belongsToHouseholdWithMandatedRestrictionRegexp    = regexp.MustCompile(`\+belongsToHousehold:\d+`)
	belongsToHouseholdWithoutMandatedRestrictionRegexp = regexp.MustCompile(`belongsToHousehold:\d+`)
)

// ensureQueryIsRestrictedToUser takes a query and userID and ensures that query
// asks that results be restricted to a given user.
func ensureQueryIsRestrictedToUser(query string, userID uint64) string {
	switch {
	case belongsToHouseholdWithMandatedRestrictionRegexp.MatchString(query):
		return query
	case belongsToHouseholdWithoutMandatedRestrictionRegexp.MatchString(query):
		query = belongsToHouseholdWithoutMandatedRestrictionRegexp.ReplaceAllString(query, fmt.Sprintf("+belongsToHousehold:%d", userID))
	case !belongsToHouseholdWithMandatedRestrictionRegexp.MatchString(query):
		query = fmt.Sprintf("%s +belongsToHousehold:%d", query, userID)
	}

	return query
}
