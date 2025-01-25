package qbr

// Data model.
type Data struct {
	Field *Field
	Value any
}

// Set adds the specified fields to the query builder's data list.
// If a field's value is nil, it is ignored and not added.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) Set(fields ...*Data) *QueryBuilder {
	// add data to query
	for _, field := range fields {
		// check is value is nil
		if isZero(field.Value) {
			continue
		}

		// add data
		qb.data = append(qb.data, *field)
	}

	// return query builder
	return qb
}

// NewData creates a new Data instance with the specified field and value.
// It initializes the Data model with the given field and value parameters,
// and returns the constructed Data object.
func NewData(field *Field, value any) *Data {
	return &Data{
		Field: field,
		Value: value,
	}
}
