package qbr

// Select adds the specified fields to the query builder's select list.
// Each field is appended to the selects slice within the QueryBuilder.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) Select(fields ...*Field) *QueryBuilder {
	// set select to null
	qb.selects = nil

	// add fields to query builder
	for _, field := range fields {
		qb.selects = append(qb.selects, *field)
	}

	// return query builder
	return qb
}

// GetSelectFields returns a copy of the query builder's select fields.
// The returned slice is a copy of the internal selects slice and is safe to modify.
func (qb *QueryBuilder) GetSelectFields() []Field {
	// conditions for returning
	fields := make([]Field, len(qb.selects))

	// copy query builder conditions
	copy(fields, qb.selects)

	// return copy conditions
	return fields
}
