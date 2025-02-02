package sqlbuilder

import (
	"fmt"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// buildConditions translates a condition slice to a SQL query string and its params.
// plc is the placeholder character to use, and params is the parameter slice to append to.
// join is the operator to use to join the condition strings, default is "AND".
// It returns the query string, the updated parameter slice, and an error if any.
func buildConditions(conds []domain.Condition, plc domain.SqlPlaceholder, params []any, join ...string) (string, []any, error) {
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
func handleLogicalCondition(cond domain.Condition, params []any, plc domain.SqlPlaceholder, lgOp domain.OperatorType) (string, []any, error) {
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
	subQuery, subParams, err := buildConditions(value, plc, params, subJoin)
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
// It takes a Condition object, a domain.SqlPlaceholder for parameter substitution, and a parameter index.
// The function checks if the condition's value is of type ValueType and handles null values accordingly.
// It retrieves the SQL operator for the given condition's operator, and constructs the SQL condition string
// with the placeholder. If the value type or operator is not supported, it returns an error.
//
// The function returns the SQL condition string, the condition's value as a parameter, and an error if any.
func handleSimpleCondition(cond domain.Condition, plc domain.SqlPlaceholder, paramIndex int) (string, any, error) {
	// check if the value type is ValueType
	if v, ok := cond.Value.(domain.ValueType); ok {
		if v == domain.ValueNull {
			// handle null value condition
			if cond.Operator == domain.OperatorNotEqual {
				return fmt.Sprintf("%s IS NOT NULL", getFieldName(cond.Field)), nil, nil
			}

			// return conditional string and success
			return fmt.Sprintf("%s IS NULL", getFieldName(cond.Field)), nil, nil
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
	condStr := fmt.Sprintf("%s %s %s", getFieldName(cond.Field), operator, getPlaceholder(plc, paramIndex))

	// return condition string, value and success
	return condStr, cond.Value, nil
}
