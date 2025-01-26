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
	// add remove zero condition s
	qb.conditions = append(
		qb.conditions,
		removeZeroCondition(conds...)...,
	)

	// return query builder
	return qb
}

// removeZeroCondition filters out conditions that have zero values from the provided
// slice of conditions. It recursively processes nested conditions and removes any
// condition with a zero value, as determined by the isZero function. The function
// returns a slice of conditions that do not contain any zero values.
func removeZeroCondition(conds ...Condition) []Condition {
	// conditions for return
	result := []Condition{}

	// stack
	stack := append([]Condition{}, conds...)

	// check all conditions
	for len(stack) > 0 {
		// condition
		cond := stack[0]
		// remove condition from stack
		stack = stack[1:]

		// select condition type
		switch v := cond.Value.(type) {
		case []Condition:
			// set new removed conditions
			cond.Value = removeZeroCondition(v...)

			// add formatted conditions
			result = append(result, cond)
		default:
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
