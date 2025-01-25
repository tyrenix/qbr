package qbr

// Condition type.
type ConditionType int

// Condition types.
const (
	ConditionEqual ConditionType = iota
	ConditionNotEqual
	ConditionLessThan
	ConditionGreaterThan
	ConditionLessThanOrEqual
	ConditionGreaterThanOrEqual
	ConditionOr
	ConditionAnd
)

// Condition model.
type Condition struct {
	Field    *Field
	Operator ConditionType
	Value    any
}

// Condition or.
//
// conds1 OR conds2 OR conds3 and so on
func Or(conds ...Condition) Condition {
	return Condition{
		Operator: ConditionOr,
		Value:    conds,
	}
}

// Condition and.
//
// conds1 AND conds2 AND conds3 and so on
func And(conds ...Condition) Condition {
	return Condition{
		Operator: ConditionAnd,
		Value:    conds,
	}
}

// Condition for equals.
//
// field = val
func Eq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionEqual,
		Value:    val,
	}
}

// Condition for not equals.
//
// field != val
func NoEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionNotEqual,
		Value:    val,
	}
}

// Condition for less than.
//
// field < val
func Lt(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionLessThan,
		Value:    val,
	}
}

// Condition for greater than.
//
// field > val
func Gt(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionGreaterThan,
		Value:    val,
	}
}

// Condition for less than or equal.
//
// field <= val
func LtOrEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionLessThanOrEqual,
		Value:    val,
	}
}

// Condition for greater than or equal.
//
// field >= val
func GtOrEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: ConditionGreaterThanOrEqual,
		Value:    val,
	}
}

// Where add filters with condition.
func (qb *QueryBuilder) Where(conds ...Condition) *QueryBuilder {
	// add filters with condition to query
	for _, cond := range conds {
		// check is value is nil
		if isZero(cond.Value) {
			continue
		}

		// add condition
		qb.conditions = append(qb.conditions, cond)
	}

	// return query builder
	return qb
}
