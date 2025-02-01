package qbr

import "github.com/tyrenix/qbr/domain"

// NewNullValue creates a new ValueType instance with the ValueNull value.
//
// It returns the domain.ValueNull constant value directly, so it can be used as a
// simple way to get the null value type.
func NewNullValue() domain.ValueType {
	return domain.ValueNull
}
