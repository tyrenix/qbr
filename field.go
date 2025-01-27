package qbr

import "reflect"

// Aggregation type.
type AggregationType int

// Aggregation types.
const (
	AggregationNone AggregationType = iota
	AggregationCount
	AggregationSum
)

// Field model.
type Field struct {
	DB          string          // DB field name.
	Aggregation AggregationType // Aggregation type.
	IgnoreOn    []OperationType // Slice with ignored operations.
}

// NewField creates a new Field model with the specified options.
//
// It takes a variable number of FieldOption functions as arguments, and applies
// each one to the created Field model before returning it. The returned Field
// model is fully constructed and ready to use.
//
// The FieldOption functions can be used to set the DB field name, ignored
// operations, and aggregation type for the created Field model.
func NewField(options ...FieldOption) *Field {
	// field
	f := &Field{}

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
func NewFieldFromStruct(s any, fieldName string) *Field {
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
type FieldOption func(*Field)

// WithDB sets the DB field for a Field model.
func WithDB(db string) FieldOption {
	return func(f *Field) {
		f.DB = db
	}
}

// WithIgnoredOperations sets the ignored operations for a Field model.
func WithIgnoredOperations(ignored ...OperationType) FieldOption {
	return func(f *Field) {
		f.IgnoreOn = append(f.IgnoreOn, ignored...)
	}
}

// WithAggregationType sets the aggregation type for a Field model.
func WithAggregationType(aggType AggregationType) FieldOption {
	return func(f *Field) {
		f.Aggregation = aggType
	}
}

// NewAllField returns a new Field model with DB type set to "*".
//
// The returned Field model is equivalent to calling NewField("*").
func NewAllField() *Field {
	return &Field{
		DB:          "*",
		Aggregation: AggregationNone,
	}
}

// NewSumField creates a new Field model with sum aggregation type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with AggregationType set to AggregationSum.
//
// Returns the created Field model with sum aggregation type.
func NewSumField(field *Field) *Field {
	return &Field{
		DB:          field.DB,
		Aggregation: AggregationSum,
	}
}

// NewCountField creates new Field model with count type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with Type set to FieldCount.
//
// Returns created Field model with count type.
func NewCountField(field *Field) *Field {
	return &Field{
		DB:          field.DB,
		Aggregation: AggregationCount,
	}
}
