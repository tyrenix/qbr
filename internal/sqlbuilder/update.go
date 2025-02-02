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
		// add data to sets
		sets = append(
			sets,
			fmt.Sprintf("%s = %s",
				getFieldName(data.Field),
				getPlaceholder(placeholder, len(params)+1),
			),
		)

		// create database value
		v, err := valueToDBValue(data.Value)
		if err != nil {
			return "", nil, err
		}

		// add params
		params = append(params, v)
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

	// return query, params and success
	return query, params, nil
}
