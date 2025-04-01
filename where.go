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

// In returns a condition that checks if the value of the given field is in the specified values.
//
// field IN (val[0], val[1], ...)
func In(field *domain.Field, val ...any) domain.Condition {
	return domain.Condition{
		Field:    field,
		Operator: domain.OperatorIn,
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
