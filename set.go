package qbr

import (
	"reflect"

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

// extractDataFromStruct extracts fields from a given struct and returns them as a slice of Data.
// If the input is a pointer, it dereferences it before processing. The function checks if the input
// is a valid struct type and iterates through its fields. For each field, it retrieves the field's
// value and annotation, and constructs a Data object. Fields with a nil value or that do not have
// a "db" annotation are ignored. The resulting slice of Data objects is returned, representing the
// struct's fields ready for inclusion in a query.
func extractDataFromStruct(s any) []*domain.Data {
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
	var data []*domain.Data

	// we go through the fields of the structure
	for i := 0; i < val.NumField(); i++ {
		// struct field
		field := val.Field(i)
		// field type
		ft := t.Field(i)

		// add data
		data = append(data, NewData(
			extractFieldFromStruct(ft),
			field.Interface(),
		))
	}

	// return query
	return data
}
