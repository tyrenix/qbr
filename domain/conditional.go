package domain

// Condition model.
type Condition struct {
	Field    *Field
	Operator OperatorType
	Value    any
}
