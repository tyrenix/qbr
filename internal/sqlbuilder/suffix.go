package sqlbuilder

import "fmt"

// buildSuffix adds a suffix to the query.
func buildSuffix(query string, suffix string) string {
	// no suffix
	if suffix == "" {
		return query
	}

	// build suffix
	return fmt.Sprintf("%s %s", query, suffix)
}
