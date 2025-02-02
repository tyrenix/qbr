package sqlbuilder

import (
	"fmt"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// CreateInsertSql creates a SQL INSERT query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func CreateInsertSql(qb Query, table string, placeholder domain.SqlPlaceholder) (string, []any, error) {
	var columns []string
	var values []string
	var params []any

	// select fields
	selects := qb.GetSelects()
	// data
	setData := qb.GetData()

	// create main query
	for _, data := range setData {
		// add database column
		columns = append(columns, getFieldName(data.Field))

		// add value
		values = append(values, getPlaceholder(placeholder, len(params)+1))

		// create database value
		v, err := valueToDBValue(data.Value)
		if err != nil {
			return "", nil, err
		}

		// add params
		params = append(params, v)
	}

	// create query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	// build returning fields
	if len(selects) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelects(selects)
	}

	// return query, params and success
	return query, params, nil
}
