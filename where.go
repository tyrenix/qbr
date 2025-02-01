package qbr

// Operator type.
type OperatorType int

// Operator types.
const (
	OperatorEqual OperatorType = iota
	OperatorNotEqual
	OperatorLessThan
	OperatorGreaterThan
	OperatorLessThanOrEqual
	OperatorGreaterThanOrEqual
	OperatorOr
	OperatorAnd
)

// Condition model.
type Condition struct {
	Field    *Field
	Operator OperatorType
	Value    any
}

// Or returns a condition that checks if any of the given conditions are true.
//
// conds1 OR conds2 OR conds3 and so on
func Or(conds ...Condition) Condition {
	return Condition{
		Operator: OperatorOr,
		Value:    conds,
	}
}

// And returns a condition that checks if all the given conditions are true.
//
// conds1 AND conds2 AND conds3 and so on
func And(conds ...Condition) Condition {
	return Condition{
		Operator: OperatorAnd,
		Value:    conds,
	}
}

// Condition for equals.
//
// Eq returns a condition that checks if the value of the given field is equal to the given value.
func Eq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorEqual,
		Value:    val,
	}
}

// Condition for not equals.
//
// NoEq returns a condition that checks if the value of the given field is not equal to the given value.
func NoEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorNotEqual,
		Value:    val,
	}
}

// Lt returns a condition that checks if the value of the given field is less than the specified value.
//
// field < val
func Lt(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorLessThan,
		Value:    val,
	}
}

// Gt returns a condition that checks if the value of the given field is greater than the specified value.
//
// field > val
func Gt(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorGreaterThan,
		Value:    val,
	}
}

// LtOrEq returns a condition that checks if the value of the given field is less than or equal to the specified value.
//
// field <= val
func LtOrEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorLessThanOrEqual,
		Value:    val,
	}
}

// GtOrEq returns a condition that checks if the value of the given field is greater than or equal to the specified value.
//
// field >= val
func GtOrEq(field *Field, val any) Condition {
	return Condition{
		Field:    field,
		Operator: OperatorGreaterThanOrEqual,
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

// GetConditions returns the conditions set for the query builder, or an empty slice if no conditions have been set.
func (qb *QueryBuilder) GetConditions() []Condition {
	// conditions for returning
	conds := make([]Condition, len(qb.conditions))

	// copy query builder conditions
	copy(conds, qb.conditions)

	// return copy conditions
	return conds
}

// removeZeroCondition removes zero conditions from the given conditions.
//
// It takes a variable number of Condition values as arguments, and returns a slice of Condition
// values with the zero conditions removed.
//
// Zero conditions are conditions that have a zero value or a null pointer value.
//
// The function also handles the case where the condition's value is a slice of Condition values,
// in which case it recursively calls itself to remove zero conditions from the sub-conditions.
//
// The function returns a new slice of Condition values with the zero conditions removed.
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
			// check is not ignored
			if isFieldIgnored(cond.Field, OperationRead) {
				continue
			}

			// check if system conditional and handle value null case
			switch t := cond.Value.(type) {
			case ValueType:
				if t == ValueNull {
					// skip if value is null and not supported aggregation or operator
					if cond.Field.Aggregation != AggregationNone ||
						(cond.Operator != OperatorEqual && cond.Operator != OperatorNotEqual) {
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
