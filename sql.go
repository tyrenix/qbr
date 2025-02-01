package qbr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// SqlPlaceholder type.
type SqlPlaceholder string

// To sql placeholders variables.
const (
	SqlDollar   SqlPlaceholder = "$"
	SqlQuestion SqlPlaceholder = "?"
)

// sqlAggregationFormats is a map that defines SQL aggregation formats for different AggregationTypes.
// It currently supports all supported aggregation types.
var sqlAggregationFormats = map[domain.AggregationType]string{
	domain.AggregationNone:  "%s",
	domain.AggregationCount: "COUNT(%s)",
	domain.AggregationSum:   "SUM(%s)",
}

// sqlOperators is a map that defines SQL operators for different OperatorTypes.
// It currently supports all supported operators.
var sqlOperators = map[domain.OperatorType]string{
	domain.OperatorAnd:                "AND",
	domain.OperatorOr:                 "OR",
	domain.OperatorEqual:              "=",
	domain.OperatorNotEqual:           "!=",
	domain.OperatorLessThan:           "<",
	domain.OperatorGreaterThan:        ">",
	domain.OperatorLessThanOrEqual:    "<=",
	domain.OperatorGreaterThanOrEqual: ">=",
}

// ToSql translates the Query into a SQL query string, its parameters, and an error if the query
// could not be built. It supports SELECT, INSERT, UPDATE, and DELETE queries. It returns an error if the
// query type is not supported.
func (qb *Query) ToSql(table string, placeholder SqlPlaceholder) (string, []any, error) {
	// select need method for build
	switch qb.operation {
	case domain.OperationRead:
		return qb.toSelectSql(table, placeholder)
	case domain.OperationCreate:
		return qb.toInsertSql(table, placeholder)
	case domain.OperationUpdate:
		return qb.toUpdateSql(table, placeholder)
	case domain.OperationDelete:
		return qb.toDeleteSql(table, placeholder)
	default:
		return "", nil, fmt.Errorf("unsupported query type: %v", qb.operation)
	}
}

// toInsertSql creates a SQL INSERT query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func (qb *Query) toInsertSql(table string, placeholder SqlPlaceholder) (string, []any, error) {
	var columns []string
	var values []string
	var params []any

	// create main query
	for _, data := range qb.data {
		// add database column
		columns = append(columns, getDBFieldName(data.Field))

		// add value
		values = append(values, buildDBPlaceholder(placeholder, len(params)+1))

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
	if len(qb.selects) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelectFields(qb.selects)
	}

	// return query, params and success
	return query, params, nil
}

// toSelectSql creates a SQL SELECT query from the Query's select list, conditions,
// sort, limit, and offset. It returns the query string, the parameters for the query,
// and an error if the query could not be built.
func (qb *Query) toSelectSql(table string, placeholder SqlPlaceholder) (string, []any, error) {
	// create main query
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		buildSelectFields(qb.selects), // create select query
		table,
	)

	// create params
	var params []any

	// is conditions exists add conditions and params
	if len(qb.conditions) > 0 {
		// create conditions
		cond, condParams, err := buildSqlConditions(qb.conditions, placeholder, nil)
		if err != nil {
			return "", nil, err
		}

		// add conditions
		query += " WHERE " + cond
		params = append(params, condParams...)
	}

	// add sort
	if len(qb.sort) > 0 {
		// order by query string
		orderByStr := " ORDER BY"

		// create order by
		for _, orderBy := range qb.sort {
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
	if v := buildSqlLimitAndOffset(qb.limit, qb.offset); v != "" {
		// add limit and offset
		query += " " + v
	}

	// return query, params and success
	return query, params, nil
}

// toUpdateSql creates a SQL UPDATE query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func (qb *Query) toUpdateSql(table string, placeholder SqlPlaceholder) (string, []any, error) {
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
				buildDBPlaceholder(placeholder, len(params)+1),
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
	if len(qb.conditions) > 0 {
		// create conditions
		conds, condsParams, err := buildSqlConditions(qb.conditions, placeholder, params)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = condsParams
	}

	// build returning fields
	if len(qb.selects) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelectFields(qb.selects)
	}

	// return query, params and success
	return query, params, nil
}

// toDeleteSql creates a SQL DELETE query from the Query's data. It returns the query string,
// the parameters for the query, and an error if the query could not be built.
func (qb *Query) toDeleteSql(table string, placeholder SqlPlaceholder) (string, []any, error) {
	var params []any

	// create base query
	query := fmt.Sprintf("DELETE FROM %s", table)

	// if exists conditions add to query
	if len(qb.conditions) > 0 {
		// create conditions
		conds, condsParams, err := buildSqlConditions(qb.conditions, placeholder, nil)
		if err != nil {
			return "", nil, err
		}

		// add conditions to query
		query += " WHERE " + conds
		// add condition params to params
		params = append(params, condsParams...)
	}

	// build returning fields
	if len(qb.selects) > 0 {
		// create returning fields
		query += " RETURNING " + buildSelectFields(qb.selects)
	}

	// return query, params and success
	return query, params, nil
}

// buildSqlLimitAndOffset creates a LIMIT and OFFSET SQL clause from the given limit and offset values.
// It returns the clause string.
func buildSqlLimitAndOffset(limit, offset uint64) string {
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

// buildSqlConditions translates a condition slice to a SQL query string and its params.
// plc is the placeholder character to use, and params is the parameter slice to append to.
// join is the operator to use to join the condition strings, default is "AND".
// It returns the query string, the updated parameter slice, and an error if any.
func buildSqlConditions(conds []domain.Condition, plc SqlPlaceholder, params []any, join ...string) (string, []any, error) {
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
		case domain.OperatorAnd, domain.OperatorOr: // for logical operator: OR, AND
			// create sub query and params
			subQuery, subParams, err := handleLogicalCondition(cond, params, plc, cond.Operator)
			if err != nil {
				return "", nil, err
			}

			// add sub query
			condStrs = append(condStrs, fmt.Sprintf("(%s)", subQuery))
			// add sub params
			params = subParams
		default: // for simple operator, >, <, <=, and so on
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

			// is param not nil add param
			if param != nil {
				params = append(params, param)
			}
		}
	}

	// create query
	query := strings.Join(condStrs, fmt.Sprintf(" %s ", opj))

	// return query, params and success
	return query, params, nil
}

// handleLogicalCondition processes a logical condition (AND/OR) within a query,
// generating a SQL sub-query and its corresponding parameters.
//
// It takes a Condition object representing the logical condition, a slice of
// current parameter values, a placeholder for SQL parameter substitution, and
// the logical operator type (AND/OR). The function validates the condition's
// value as a slice of sub-conditions, then recursively builds SQL sub-queries
// for each condition within the logical group. The resulting SQL string and
// updated parameter list are returned, along with an error if any occurs
// during the process.
func handleLogicalCondition(cond domain.Condition, params []any, plc SqlPlaceholder, lgOp domain.OperatorType) (string, []any, error) {
	// assert type
	value, ok := cond.Value.([]domain.Condition)
	if !ok {
		return "", nil, fmt.Errorf("invalid value for logical operator %d", lgOp)
	}

	// sub queries
	var subQueries []string

	// sub join
	subJoin := "AND"
	if lgOp == domain.OperatorOr {
		subJoin = "OR"
	}

	// create sub query
	subQuery, subParams, err := buildSqlConditions(value, plc, params, subJoin)
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

// handleSimpleCondition processes a simple condition within a SQL query, generating a SQL condition string
// and its corresponding parameter.
//
// It takes a Condition object, a SqlPlaceholder for parameter substitution, and a parameter index.
// The function checks if the condition's value is of type ValueType and handles null values accordingly.
// It retrieves the SQL operator for the given condition's operator, and constructs the SQL condition string
// with the placeholder. If the value type or operator is not supported, it returns an error.
//
// The function returns the SQL condition string, the condition's value as a parameter, and an error if any.
func handleSimpleCondition(cond domain.Condition, plc SqlPlaceholder, paramIndex int) (string, any, error) {
	// check if the value type is ValueType
	if v, ok := cond.Value.(domain.ValueType); ok {
		if v == domain.ValueNull {
			// handle null value condition
			if cond.Operator == domain.OperatorNotEqual {
				return fmt.Sprintf("%s IS NOT NULL", getDBFieldName(cond.Field)), nil, nil
			}

			// return conditional string and success
			return fmt.Sprintf("%s IS NULL", getDBFieldName(cond.Field)), nil, nil
		}

		// return error
		return "", nil, fmt.Errorf("unsupported value type: %d", v)
	}

	// get SQL operator
	operator := getSqlOperator(cond.Operator)
	if operator == "" {
		return "", nil, fmt.Errorf("unsupported operator: %d", cond.Operator)
	}

	// create condition string with placeholder
	condStr := fmt.Sprintf("%s %s %s", getDBFieldName(cond.Field), operator, buildDBPlaceholder(plc, paramIndex))

	// return condition string, value and success
	return condStr, cond.Value, nil
}

// getSqlOperator returns the SQL operator associated with the given OperatorType.
// If the operator is not supported, it returns an empty string.
func getSqlOperator(op domain.OperatorType) string {
	// get operator
	v, ok := sqlOperators[op]
	if !ok {
		return ""
	}

	// return operator
	return v
}

// buildSelectFields formats a slice of Field objects into a SQL select statement string.
// It iterates over the provided fields, and for each field, it checks if there is an
// associated SQL format in the sqlFieldFormats map based on the field's type. If a format
// exists, it retrieves the database field name and applies the format, adding the result
// to the list of select fields. The function returns a comma-separated string of the
// formatted select fields.
func buildSelectFields(fields []domain.Field) string {
	// fields
	var result []string

	// iterate over the slice of fields
	for _, field := range fields {
		// check is contains in map
		format, ok := sqlAggregationFormats[field.Aggregation]
		if !ok {
			continue
		}

		// append the formatted field to the result slice
		result = append(
			result,
			fmt.Sprintf(
				format,
				getDBFieldName(&field), // get database field name
			), // create field name
		)
	}

	// return the fields as a comma-separated string
	return strings.Join(result, ", ")
}

// getDBFieldName takes a Field object and returns the string value of its DB
// field. This is the field name in the database that the field corresponds to.
func getDBFieldName(field *domain.Field) string {
	return string(field.DB)
}

// buildPlaceholder generates a SQL placeholder string based on the specified
// placeholder type and index. If the placeholder type is ToSqlDollar, it returns
// a parameterized string using the dollar sign format (e.g., $1, $2). Otherwise,
// it returns the placeholder type as a string.
func buildDBPlaceholder(plc SqlPlaceholder, index int) string {
	// dollar placeholder
	if plc == SqlDollar {
		return fmt.Sprintf("$%d", index)
	}

	// return default placeholder
	return string(plc)
}

// valueToDBValue takes a value and returns a value that can be used in a SQL query
// and an error if the value is not supported.
//
// If the value is of type ValueType, it returns the string "NULL" if the value
// is ValueNull, otherwise it returns an error.
//
// If the value is a struct or a pointer to a struct, it converts the value to
// a JSON string and returns it.
//
// If the value is not a struct or a pointer to a struct, it returns the original
// value.
func valueToDBValue(value any) (any, error) {
	// is value is ValueType
	if v, ok := value.(domain.ValueType); ok {
		// is null value return null
		if v == domain.ValueNull {
			return nil, nil
		}

		// return nil
		return nil, fmt.Errorf("unsupported value type: %d", v)
	}

	// get reflect value
	v := reflect.ValueOf(value)

	// check if the value is a struct or a pointer to a struct
	if v.Kind() == reflect.Ptr &&
		v.Elem().Kind() == reflect.Struct {
		// convert struct to JSON
		j, err := json.Marshal(v.Elem().Interface())
		if err != nil {
			return nil, err
		}

		// return json string
		return j, nil
	}

	// return the original value if not a struct or pointer to a struct
	return value, nil
}
