package sqlbuilder

import (
	"fmt"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// CreateSelectSql creates a SQL SELECT query from the Query's select list, conditions,
// sort, limit, and offset. It returns the query string, the parameters for the query,
// and an error if the query could not be built.
func CreateSelectSql(qb Query, table string, placeholder domain.SqlPlaceholder) (string, []any, error) {
	// create main query
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		buildSelects(qb.GetSelects()), // create select query
		table,
	)

	// create params
	var params []any

	// conditionals
	conds := qb.GetConditions()
	// sorts
	sorts := qb.GetSort()
	// limit
	limit := qb.GetLimit()
	// offset
	offset := qb.GetOffset()

	// is conditions exists add conditions and params
	if len(conds) > 0 {
		// create conditions
		cond, condParams, err := buildConditions(conds, placeholder, nil)
		if err != nil {
			return "", nil, err
		}

		// add conditions
		query += " WHERE " + cond
		params = append(params, condParams...)
	}

	// add sort
	if len(sorts) > 0 {
		// order by query string
		orderByStr := " ORDER BY"

		// create order by
		sortClauses := make([]string, len(sorts))
		for i, sort := range sorts {
			sortClauses[i] = fmt.Sprintf(
				"%s %s",
				getFieldName(sort.Field),
				sort.Type,
			)
		}

		// add order by
		query += orderByStr + " " + strings.Join(sortClauses, ", ")
	}

	// add limit and offset
	if v := buildLimitAndOffset(limit, offset); v != "" {
		// add limit and offset
		query += " " + v
	}

	// return query, params and success
	return query, params, nil
}
