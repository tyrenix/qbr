package qbr

import (
	"reflect"
	"strings"
)

// Data model.
type Data struct {
	Field *Field
	Value any
}

// Set adds the specified fields to the query builder's data list.
// If a field's value is nil, it is ignored and not added.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) Set(data ...*Data) *QueryBuilder {
	// add data to query
	for _, d := range data {
		// check is value is nil
		if d.Field == nil || isZero(d.Value) {
			continue
		}

		// check is ignore
		for _, i := range d.Field.Ignore {
			// this ignored, skip
			if i == qb.queryType {
				continue
			}
		}

		// add data
		qb.data = append(qb.data, *d)
	}

	// return query builder
	return qb
}

// SetStruct adds the specified struct's fields to the query builder's data list.
// If the struct is a pointer, it is dereferenced first.
// If the struct contains a field with a tag "db" and the field is not nil, it is added.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) SetStruct(s any) *QueryBuilder {
	// struct value
	val := reflect.ValueOf(s)
	// struct type
	t := reflect.TypeOf(s)

	// check is pointer
	if val.Kind() == reflect.Ptr {
		// check is nil
		if val.IsNil() {
			return qb
		}

		// dereference pointer
		val = val.Elem()
		t = t.Elem()
	}

	// check is struct
	if val.Kind() != reflect.Struct {
		return qb
	}

	// we go through the fields of the structure
	for i := 0; i < val.NumField(); i++ {
		// struct field
		field := val.Field(i)
		// field type
		ft := t.Field(i)

		// set data
		qb.Set(NewData(
			extractFieldFromStructAnnotation(ft),
			field.Interface(),
		))
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

// extractFieldFromStructAnnotation takes a reflect.StructField and extracts the
// "db" tag from the field's struct tag. If the tag is not present, it returns nil.
// Otherwise, it creates a new Field instance with the extracted tag value and
// returns it.
func extractFieldFromStructAnnotation(ft reflect.StructField) *Field {
	// get tags from field annotation
	db := ft.Tag.Get(string(QueryBuilderDB))

	// check is not empty
	if db == "" {
		return nil
	}

	// ignores
	var ignores []QueryBuilderType

	// extract ignore
	for _, ignore := range strings.Split(ft.Tag.Get(string(QueryBuilderIgnore)), ",") {
		// check is empty
		if ignore == "" {
			continue
		}

		// add ignore
		ignores = append(
			ignores,
			QueryBuilderType(ignore),
		)
	}

	// create field
	field := NewField(
		FieldDB(db),
	)

	// add ignores
	field.Ignore = ignores

	// return fields
	return field
}
