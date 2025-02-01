package domain

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
