package qbr

// Select adds the specified fields to the query builder's select list.
// Each field is appended to the selects slice within the QueryBuilder.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) Select(fields ...*Field) *QueryBuilder {
	// add fields to query builder
	for _, field := range fields {
		qb.selects = append(qb.selects, *field)
	}

	// return query builder
	return qb
}

// GetSelectFields returns the select list of the query builder, which is a slice of Field models.
func (qb *QueryBuilder) GetSelectFields() []Field {
	return qb.selects
}
