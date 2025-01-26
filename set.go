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
		skip := false
		for _, i := range d.Field.Ignore {
			// this ignored, skip
			if i == qb.queryType {
				skip = true
				break
			}
		}

		// if skip
		if skip {
			continue
		}

		// add data
		qb.data = append(qb.data, *d)
	}

	// return query builder
	return qb
}

// SetStruct adds the specified struct fields to the query builder's data list.
// If a field's value is nil, it is ignored and not added.
// Returns the modified QueryBuilder instance for method chaining.
func (qb *QueryBuilder) SetStruct(s any) *QueryBuilder {
	// extract data from struct
	data := extractDataFromStruct(s)
	// set data to query builder
	return qb.Set(data...)
}

// extractDataFromStruct extracts fields from a given struct and returns them as a slice of Data.
// If the input is a pointer, it dereferences it before processing. The function checks if the input
// is a valid struct type and iterates through its fields. For each field, it retrieves the field's
// value and annotation, and constructs a Data object. Fields with a nil value or that do not have
// a "db" annotation are ignored. The resulting slice of Data objects is returned, representing the
// struct's fields ready for inclusion in a query builder.
func extractDataFromStruct(s any) []*Data {
	// struct value
	val := reflect.ValueOf(s)
	// struct type
	t := reflect.TypeOf(s)

	// check is pointer
	if val.Kind() == reflect.Ptr {
		// check is nil
		if val.IsNil() {
			return nil
		}

		// dereference pointer
		val = val.Elem()
		t = t.Elem()
	}

	// check is struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	// create data slice
	var data []*Data

	// we go through the fields of the structure
	for i := 0; i < val.NumField(); i++ {
		// struct field
		field := val.Field(i)
		// field type
		ft := t.Field(i)

		// add data
		data = append(data, NewData(
			extractFieldFromStructAnnotation(ft),
			field.Interface(),
		))
	}

	// return query builder
	return data
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
