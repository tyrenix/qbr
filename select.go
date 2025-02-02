package qbr

import "github.com/tyrenix/qbr/domain"

// Select sets the fields to be selected in the query. If no fields are
// specified, all fields are selected. The fields parameter is a variable
// argument list, so you can pass in any number of fields or an array/slice
// of fields. The method returns the QueryBuilder instance to support method
// chaining.
func (qb *Query) Select(fields ...*domain.Field) *Query {
	// set select to null
	qb.selects = nil

	// add fields to query
	for _, field := range fields {
		qb.selects = append(qb.selects, *field)
	}

	// return query
	return qb
}

// GetSelects returns the select fields set for the query builder, or an empty slice if no select fields have been set.
func (qb *Query) GetSelects() []domain.Field {
	// conditions for returning
	fields := make([]domain.Field, len(qb.selects))

	// copy query  conditions
	copy(fields, qb.selects)

	// return copy conditions
	return fields
}
