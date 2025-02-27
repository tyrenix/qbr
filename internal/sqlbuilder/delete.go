package sqlbuilder

import (
	"fmt"

	"github.com/tyrenix/qbr/domain"
)

// CreateDeleteSql creates a SQL DELETE query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func CreateDeleteSql(qb Query, table string, placeholder domain.SqlPlaceholder) (string, []any, error) {
	var params []any

	// create base query
	query := fmt.Sprintf("DELETE FROM %s", table)

	// conditionals
	conds := qb.GetConditions()

	// if exists conditions add to query
	if len(conds) > 0 {
		// create conditions
		conds, condsParams, err := buildConditions(conds, placeholder, nil)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = append(params, condsParams...)
	}

	// build returning fields
	if len(conds) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelects(qb.GetSelects())
	}

	// add suffix
	query = buildSuffix(query, qb.GetSuffix())

	// return query, params and success
	return query, params, nil
}
