package sqlbuilder

import (
	"fmt"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// CreateUpdateSql creates a SQL UPDATE query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func CreateUpdateSql(qb Query, table string, placeholder domain.SqlPlaceholder) (string, []any, error) {
	var sets []string
	var params []any

	// create base query
	query := fmt.Sprintf("UPDATE %s SET ", table)

	// select fields
	selects := qb.GetSelects()
	// conditionals
	conds := qb.GetConditions()
	// data
	setData := qb.GetData()

	// create add update params
	for _, data := range setData {
		// create set data
		field, plc, value, err := buildSetData(data, placeholder, len(params)+1)
		if err != nil {
			return "", nil, err
		}

		// add data to sets
		sets = append(sets, fmt.Sprintf("%s = %s", field, plc))
		// add value to params
		params = append(params, value)
	}

	// add to query set data
	query += strings.Join(sets, ", ")

	// if exists conditions add to query
	if len(conds) > 0 {
		// create conditions
		conds, condsParams, err := buildConditions(conds, placeholder, params)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = condsParams
	}

	// build returning fields
	if len(selects) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelects(selects)
	}

	// add suffix
	query = buildSuffix(query, qb.GetSuffix())

	// return query, params and success
	return query, params, nil
}
