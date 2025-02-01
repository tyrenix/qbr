package qbr

import "github.com/tyrenix/qbr/domain"

// Or returns a condition that checks if any of the given conditions are true.
//
// conds1 OR conds2 OR conds3 and so on
func Or(conds ...domain.Condition) domain.Condition {
	return domain.Condition{
		Operator: domain.OperatorOr,
		Value:    conds,
	}
}

// And returns a condition that checks if all the given conditions are true.
//
// conds1 AND conds2 AND conds3 and so on
func And(conds ...domain.Condition) domain.Condition {
	return domain.Condition{
		Operator: domain.OperatorAnd,
		Value:    conds,
	}
}

// Condition for equals.
//
// Eq returns a condition that checks if the value of the given field is equal to the given value.
func Eq(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorEqual,
		Value:    val,
	}
}

// Condition for not equals.
//
// NoEq returns a condition that checks if the value of the given field is not equal to the given value.
func NoEq(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorNotEqual,
		Value:    val,
	}
}

// Lt returns a condition that checks if the value of the given field is less than the specified value.
//
// field < val
func Lt(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorLessThan,
		Value:    val,
	}
}

// Gt returns a condition that checks if the value of the given field is greater than the specified value.
//
// field > val
func Gt(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorGreaterThan,
		Value:    val,
	}
}

// LtOrEq returns a condition that checks if the value of the given field is less than or equal to the specified value.
//
// field <= val
func LtOrEq(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorLessThanOrEqual,
		Value:    val,
	}
}

// GtOrEq returns a condition that checks if the value of the given field is greater than or equal to the specified value.
//
// field >= val
func GtOrEq(field *domain.Field, val any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorGreaterThanOrEqual,
		Value:    val,
	}
}

// Where adds the specified conditions to the QueryBuilder's conditions list.
// If a condition's Value is nil or zero, it is ignored and not added.
// Additionally, if the condition's Field is ignored for the current query type, it is also ignored and not added.
// The method returns the modified QueryBuilder instance for method chaining.
func (qb *Query) Where(conds ...domain.Condition) *Query {
	// add remove zero condition s
	qb.conditions = append(
		qb.conditions,
		removeZeroCondition(conds...)...,
	)

	// return query
	return qb
}

// GetConditions returns the conditions set for the query builder, or an empty slice if no conditions have been set.
func (qb *Query) GetConditions() []domain.Condition {
	// conditions for returning
	conds := make([]domain.Condition, len(qb.conditions))

	// copy query builder conditions
	copy(conds, qb.conditions)

	// return copy conditions
	return conds
}

// removeZeroCondition takes a variable number of conditions and returns a new slice
// with the following changes:
//  1. Conditions with a Value of nil or a zero value are removed.
//  2. Conditions with a Field that is ignored for the current query type are removed.
//  3. Conditions with a Value of domain.ValueNull are removed if the condition is not
//     an aggregation or an equality/inequality check.
//
// The method returns the modified slice of conditions.
func removeZeroCondition(conds ...domain.Condition) []domain.Condition {
	// conditions for return
	result := []domain.Condition{}

	// stack
	stack := append([]domain.Condition{}, conds...)

	// check all conditions
	for len(stack) > 0 {
		// condition
		cond := stack[0]
		// remove condition from stack
		stack = stack[1:]

		// select condition type
		switch v := cond.Value.(type) {
		case []domain.Condition:
			// set new removed conditions
			cond.Value = removeZeroCondition(v...)

			// add formatted conditions
			result = append(result, cond)
		default:
			// check is not ignored
			if isFieldIgnored(cond.Field, domain.OperationRead) {
				continue
			}

			// check if system conditional and handle value null case
			switch t := cond.Value.(type) {
			case domain.ValueType:
				if t == domain.ValueNull {
					// skip if value is null and not supported aggregation or operator
					if cond.Field.Aggregation != domain.AggregationNone ||
						(cond.Operator != domain.OperatorEqual && cond.Operator != domain.OperatorNotEqual) {
						continue
					}
				}
			}

			// check is not zero
			if !isZero(v) {
				// add condition
				result = append(result, cond)
			}
		}
	}

	// return conditions
	return result
}
