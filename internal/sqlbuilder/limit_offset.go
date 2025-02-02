package sqlbuilder

import (
	"fmt"
	"strings"
)

// buildLimitAndOffset creates a LIMIT and OFFSET SQL clause from the given limit and offset values.
// It returns the clause string.
func buildLimitAndOffset(limit, offset uint64) string {
	// sql query
	query := ""

	// add limit
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	// add offset
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	// create limit and offset
	return strings.TrimSpace(query)
}
