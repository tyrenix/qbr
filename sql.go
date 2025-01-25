package qbr

import (
	"fmt"
	"strings"
)

// To sql placeholder.
type ToSqlPlaceholder string

// To sql placeholders variables.
const (
	ToSqlDollar   ToSqlPlaceholder = "$"
	ToSqlQuestion ToSqlPlaceholder = "?"
)

// ToSql create sql query.
func (qb *QueryBuilder) ToSql(table string, placeholder ToSqlPlaceholder) (string, []any, error) {
	// select need method for build
	switch qb.queryType {
	case QueryBuilderRead:
		return qb.toSelectSql(table, placeholder)
	case QueryBuilderCreate:
		return qb.toInsertSql(table, placeholder)
	case QueryBuilderUpdate:
		return qb.toUpdateSql(table, placeholder)
	case QueryBuilderDelete:
		return qb.toDeleteSql(table, placeholder)
	default:
		return "", nil, fmt.Errorf("unsupported query type: %v", qb.queryType)
	}
}

// toInsertSql creates a SQL INSERT query from the QueryBuilder's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func (qb *QueryBuilder) toInsertSql(table string, placeholder ToSqlPlaceholder) (string, []any, error) {
	var columns []string
	var values []string
	var params []any

	// create main query
	for _, data := range qb.data {
		// add database column
		columns = append(columns, getDBFieldName(data.Field))

		// add value
		values = append(values, createPlaceholder(placeholder, len(params)+1))

		// add params
		params = append(params, data.Value)
	}

	// create query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	// return query, params and success
	return query, params, nil
}

// toSelectSql creates a SQL SELECT query from the QueryBuilder's select list, conditions,
// sort, limit, and offset. It returns the query string, the parameters for the query,
// and an error if the query could not be built.
func (qb *QueryBuilder) toSelectSql(table string, placeholder ToSqlPlaceholder) (string, []any, error) {
	// create main query
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		formatSelectFields(qb.selects), // create select query
		table,
	)

	// create params
	var params []any

	// add conditions
	cond, condParams, err := toSqlConditions(qb.conditions, placeholder, nil)
	if err != nil {
		return "", nil, err
	}

	// add conditions and params
	query += " WHERE " + cond
	params = append(params, condParams...)

	// add sort
	if len(qb.orderBy) > 0 {
		// order by query string
		orderByStr := " ORDER BY"

		// create order by
		for _, orderBy := range qb.orderBy {
			orderByStr += fmt.Sprintf(
				" %s %s",
				getDBFieldName(orderBy.Field),
				orderBy.Type,
			)
		}

		// add order by
		query += orderByStr
	}

	// add limit and offset
	if v := toSqlLimitAndOffset(qb.limit, qb.offset); v != "" {
		// add limit and offset
		query += " " + v
	}

	// return query, params and success
	return query, params, nil
}

func (qb *QueryBuilder) toUpdateSql(table string, placeholder ToSqlPlaceholder) (string, []any, error) {
	var sets []string
	var params []any

	// create base query
	query := fmt.Sprintf("UPDATE %s SET ", table)

	// create add update params
	for _, data := range qb.data {
		// add data to sets
		sets = append(
			sets,
			fmt.Sprintf("%s = %s",
				getDBFieldName(data.Field),
				createPlaceholder(placeholder, len(params)+1),
			),
		)

		// add params
		params = append(params, data.Value)
	}

	// add to query set data
	query += strings.Join(sets, ", ")

	// if exists conditions add to query
	if len(qb.conditions) > 0 {
		// create conditions
		conds, condsParams, err := toSqlConditions(qb.conditions, placeholder, params)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = condsParams
	}

	// add limit and offset
	if v := toSqlLimitAndOffset(qb.limit, qb.offset); v != "" {
		// add limit and offset
		query += " " + v
	}

	// return query, params and success
	return query, params, nil
}

// toDeleteSql creates a SQL DELETE query from the QueryBuilder's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func (qb *QueryBuilder) toDeleteSql(table string, placeholder ToSqlPlaceholder) (string, []any, error) {
	var params []any

	// create base query
	query := fmt.Sprintf("DELETE FROM %s", table)

	// if exists conditions add to query
	if len(qb.conditions) > 0 {
		// create conditions
		conds, condsParams, err := toSqlConditions(qb.conditions, placeholder, nil)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = append(params, condsParams...)
	}

	// add limit and offset
	if v := toSqlLimitAndOffset(qb.limit, qb.offset); v != "" {
		// add limit and offset
		query += " " + v
	}

	// return query, params and success
	return query, params, nil
}

func toSqlLimitAndOffset(limit, offset uint64) string {
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

// toSqlConditions translates a condition slice to a SQL query string and its params.
// plc is the placeholder character to use, and params is the parameter slice to append to.
// join is the operator to use to join the condition strings, default is "AND".
// It returns the query string, the updated parameter slice, and an error if any.
func toSqlConditions(conds []Condition, plc ToSqlPlaceholder, params []any, join ...string) (string, []any, error) {
	// check check conditions count
	if len(conds) == 0 {
		return "", nil, nil
	}

	// condition strings
	var condStrs []string

	// operator join
	opj := "AND"
	if len(join) > 0 {
		opj = join[0]
	}

	// condition join
	for _, cond := range conds {
		switch cond.Operator {
		case ConditionAnd, ConditionOr: // for logical condition: OR, AND
			// create sub query and params
			subQuery, subParams, err := handleLogicalCondition(cond, params, plc, cond.Operator)
			if err != nil {
				return "", nil, err
			}

			// add sub query
			condStrs = append(condStrs, fmt.Sprintf("(%s)", subQuery))
			// add sub params
			params = subParams
		default: // for simple condition, >, <, <=, and so on
			// create condition
			conditionStr, param, err := handleSimpleCondition(cond, plc, len(params)+1)
			if err != nil {
				return "", nil, err
			}

			// check is condition is empty
			if conditionStr == "" {
				continue
			}

			// add condition
			condStrs = append(condStrs, conditionStr)
			// add param
			params = append(params, param)
		}
	}

	// create query
	query := strings.Join(condStrs, fmt.Sprintf(" %s ", opj))

	// return query, params and success
	return query, params, nil
}

// handleLogicalCondition handle logical condition, ConditionAnd or ConditionOr.
// It will translate the condition to sql query and params.
// The condition value is a 2D array, each sub array is a condition group.
// The sub arrays will be joined with the logical operator, ConditionAnd or ConditionOr.
// The condition groups will be joined with the logical operator too.
func handleLogicalCondition(cond Condition, params []any, plc ToSqlPlaceholder, lgOp ConditionType) (string, []any, error) {
	// assert type
	value, ok := cond.Value.([]Condition)
	if !ok {
		return "", nil, fmt.Errorf("invalid value for logical operator %d", lgOp)
	}

	// sub queries
	var subQueries []string

	// sub join
	subJoin := "AND"
	if lgOp == ConditionOr {
		subJoin = "OR"
	}

	// create sub query
	subQuery, subParams, err := toSqlConditions(value, plc, params, subJoin)
	if err != nil {
		return "", nil, err
	}

	// append the sub-query
	subQueries = append(subQueries, subQuery)
	// append the sub-params
	params = subParams

	// return query, params and success
	return strings.Join(subQueries, fmt.Sprintf(" %s ", subJoin)), params, nil
}

// handleSimpleCondition creates a simple condition for a given condition object.
// plc is the placeholder character to use, and paramIndex is the index of the parameter
// in the parameter slice. It returns the condition string, the parameter value, and an error.
// If the operator is unsupported, it returns an error.
func handleSimpleCondition(cond Condition, plc ToSqlPlaceholder, paramIndex int) (string, any, error) {
	// get sql operator
	operator := getSqlOperator(cond.Operator)
	if operator == "" {
		return "", nil, fmt.Errorf("unsupported operator: %d", cond.Operator)
	}

	// check is value is zero
	if isZero(cond.Value) {
		return "", nil, nil
	}

	// create placeholder
	placeholder := createPlaceholder(plc, paramIndex)
	// create condition
	conditionStr := fmt.Sprintf("%s %s %s", getDBFieldName(cond.Field), operator, placeholder)

	// return condition string, val and success
	return conditionStr, cond.Value, nil
}

// getSqlOperator translates a ConditionType into a SQL operator string.
// Returns an empty string if the ConditionType is not supported.
func getSqlOperator(op ConditionType) string {
	switch op {
	case ConditionEqual:
		return "="
	case ConditionNotEqual:
		return "!="
	case ConditionLessThan:
		return "<"
	case ConditionGreaterThan:
		return ">"
	case ConditionLessThanOrEqual:
		return "<="
	case ConditionGreaterThanOrEqual:
		return ">="
	default:
		return ""
	}
}

// formatSelectFields formats a slice of Field objects into a SQL select string.
// It generates SQL expressions based on the Field type, such as SUM or COUNT
// for aggregation fields, and returns the fields as a comma-separated string.
//
// The function iterates over the slice of fields and formats each field based on
// its type. Non-aggregation fields are formatted as a simple field name. Aggregation
// fields are formatted as a SQL expression.
func formatSelectFields(fields []Field) string {
	// fields
	var result []string

	// iterate over the slice of fields
	for _, field := range fields {
		fname := getDBFieldName(&field)

		// format the field based on its type
		switch field.Type {
		case FieldSum:
			// format the field as a SUM expression
			fname = fmt.Sprintf("SUM(%s)", fname)
		case FieldCount:
			// format the field as a COUNT expression
			fname = fmt.Sprintf("COUNT(%s)", fname)
		}

		// append the formatted field to the result slice
		result = append(result, fname)
	}

	// return the fields as a comma-separated string
	return strings.Join(result, ", ")
}

// getDBFieldName takes a Field object and returns the string value of its DB
// field. This is the field name in the database that the field corresponds to.
func getDBFieldName(field *Field) string {
	return string(field.DB)
}

// createPlaceholder generates a SQL placeholder string based on the specified
// placeholder type and index. If the placeholder type is ToSqlDollar, it returns
// a parameterized string using the dollar sign format (e.g., $1, $2). Otherwise,
// it returns the placeholder type as a string.
func createPlaceholder(plc ToSqlPlaceholder, index int) string {
	// dollar placeholder
	if plc == ToSqlDollar {
		return fmt.Sprintf("$%d", index)
	}

	// return default placeholder
	return string(plc)
}
