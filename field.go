package qbr

import (
	"reflect"

	"github.com/tyrenix/qbr/domain"
)

// NewField creates a new Field model with the specified options.
//
// It takes a variable number of FieldOption functions as arguments, and applies
// each one to the created Field model before returning it. The returned Field
// model is fully constructed and ready to use.
//
// The FieldOption functions can be used to set the DB field name, ignored
// operations, and aggregation type for the created Field model.
func NewField(options ...FieldOption) *domain.Field {
	// field
	f := &domain.Field{}

	// add all options to field
	for _, opt := range options {
		opt(f)
	}

	// return field
	return f
}

// NewFieldFromStruct creates a Field model from a specified field in a struct.
// It takes an input struct 's' and a 'fieldName' string, and returns a pointer
// to a Field model representing the specified field. If 's' is a pointer, it is
// dereferenced before processing. If 's' is not a struct, or if the specified
// field does not exist, the function returns nil.
func NewFieldFromStruct(s any, fieldName string) *domain.Field {
	val := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// if pointer, dereference
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		t = t.Elem()
	}

	// check is struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	// find field by name
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return nil
	}

	// extract field from struct
	return extractFieldFromStruct(field)
}

// FieldOption is a function that configures a Field model.
type FieldOption func(*domain.Field)

// WithDB sets the DB field for a Field model.
func WithDB(db string) FieldOption {
	return func(f *domain.Field) {
		f.DB = db
	}
}

// WithIgnoreOn returns a FieldOption that sets the ignored operations for a Field model.
//
// It takes a variable number of OperationType values as arguments, and returns a FieldOption that
// appends each argument to the IgnoreOn field of the Field model. The returned FieldOption can be
// used to configure a Field model created by NewField.
func WithIgnoreOn(ignoreOn ...domain.OperationType) FieldOption {
	return func(f *domain.Field) {
		f.IgnoreOn = append(f.IgnoreOn, ignoreOn...)
	}
}

// WithAggregation sets the aggregation type for a Field model.
//
// It takes an AggregationType and returns a FieldOption that sets the
// Aggregation field of the Field model to the specified value.
func WithAggregation(agg domain.AggregationType) FieldOption {
	return func(f *domain.Field) {
		f.Aggregation = agg
	}
}

// NewAllField returns a new Field model with DB type set to "*".
//
// The returned Field model is equivalent to calling NewField("*").
func NewAllField() *domain.Field {
	return &domain.Field{
		DB:          "*",
		Aggregation: domain.AggregationNone,
	}
}

// NewSumField creates a new Field model with sum aggregation type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with AggregationType set to AggregationSum.
//
// Returns the created Field model with sum aggregation type.
func NewSumField(field *domain.Field) *domain.Field {
	return &domain.Field{
		DB:          field.DB,
		Aggregation: domain.AggregationSum,
	}
}

// NewCountField creates new Field model with count type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with Type set to FieldCount.
//
// Returns created Field model with count type.
func NewCountField(field *domain.Field) *domain.Field {
	return &domain.Field{
		DB:          field.DB,
		Aggregation: domain.AggregationCount,
	}
}

// IsFieldEqual checks if two Field objects are equal by comparing their
// DB field names. If either of the input Field objects is nil, the function
// returns false.
func IsFieldEqual(field1, field2 *domain.Field) bool {
	// is fields is nil
	if field1 == nil || field2 == nil {
		return false
	}

	// check is field equals
	return field1.DB == field2.DB
}
