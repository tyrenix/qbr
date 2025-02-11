package qbr

import (
	"github.com/tyrenix/qbr/domain"
)

// NewData creates a new instance of domain.Data with the specified field and value.
// It takes a pointer to a domain.Field and a value of any type as arguments.
// Returns a pointer to the newly created domain.Data object containing the provided field and value.
func NewData(field *domain.Field, value any) *domain.Data {
	return &domain.Data{
		Field: field,
		Value: value,
	}
}

// Set adds the specified Data objects to the QueryBuilder's data list. If a Data object's Value is
// nil or zero, it is ignored and not added. Additionally, if the Data object's Field is ignored for
// the current query type, it is also ignored and not added. Returns the modified QueryBuilder
// instance for method chaining.
func (qb *Query) Set(data ...*domain.Data) *Query {
	// add data to query
	for _, d := range data {
		// check is value is nil
		if d.Field == nil || isZero(d.Value) {
			continue
		}

		// check is ignore
		if isFieldIgnored(d.Field, qb.operation) {
			continue
		}

		// add data
		qb.data = append(qb.data, *d)
	}

	// return query
	return qb
}

// SetStruct adds the fields of the given struct to the QueryBuilder's data list, excluding any fields
// with a nil value or that do not have a "db" annotation. The struct is first dereferenced if it is a
// pointer. The method returns the modified QueryBuilder instance for method chaining.
func (qb *Query) SetStruct(s any) *Query {
	// extract data from struct
	data := extractDataFromStruct(s)

	// set data to query
	return qb.Set(data...)
}

// GetData returns the data set for the query, or an empty slice if no data has been set.
func (qb *Query) GetData() []domain.Data {
	// init new data slice
	data := make([]domain.Data, len(qb.data))

	// copy slice
	copy(data, qb.data)

	// return data
	return data
}

// Add creates a Modification model that applies the addition operator
// to the value of the specified field. It takes a pointer to a domain.Field and
// a value of any type as arguments. The returned Modification model is
// configured to add the given value to the field's value.
func Add(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationAdd,
	}
}

// Subtract creates a Modification model that applies the subtraction operator
// to the value of the specified field. It takes a pointer to a domain.Field and
// a value of any type as arguments. The returned Modification model is
// configured to subtract the given value from the field's value.
func Subtract(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationSubtract,
	}
}

// Multiply creates a Modification model that applies the multiplication operator
// to the value of the specified field. It takes a pointer to a domain.Field and
// a value of any type as arguments and returns a pointer to the created
// domain.Modification object with the field, multiplication operator, and
// value.
func Multiply(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationMultiply,
	}
}

// Divide creates a Modification model that applies the division operator
// to the value of the specified field. It takes a pointer to a domain.Field and
// a value of any type as arguments and returns a pointer to the created
// domain.Modification object with the field, division operator, and value.
func Divide(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationDivide,
	}
}

// BitwiseAnd creates a Modification model that applies the bitwise AND
// operator to the value of the specified field. It takes a pointer to a
// domain.Field and a value of any type as arguments and returns a pointer to
// the created domain.Modification object with the field, AND operator, and
// value.
func BitwiseAnd(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationBitwiseAnd,
	}
}

// BitwiseOr creates a Modification model that applies the bitwise OR
// operator to the value of the specified field. It takes a pointer to a
// domain.Field and a value of any type as arguments and returns a pointer to
// the created domain.Modification object with the field, OR operator, and
// value.
func BitwiseOr(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationBitwiseOr,
	}
}

// BitwiseXor creates a Modification model that applies the bitwise XOR
// operator to the value of the specified field. It takes a pointer to a
// domain.Field and a value of any type as arguments and returns a pointer to
// the created domain.Modification object with the field, XOR operator, and
// value.
func BitwiseXor(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationBitwiseXor,
	}
}

// ShiftLeft returns a Modification model with the given field and value.
// It applies the bitwise shift left operator to the value of the specified
// field.
func ShiftLeft(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationShiftLeft,
	}
}

// ShiftRight returns a Modification model with the given field and value.
// It applies the bitwise shift right operator to the value of the specified
// field.
func ShiftRight(field *domain.Field, value any) *domain.Modification {
	return &domain.Modification{
		Field:    field,
		Value:    value,
		Operator: domain.ModificationShiftRight,
	}
}
