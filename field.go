package qbr

// Field type.
type FieldType int

const (
	FieldDefault FieldType = iota
	FieldAll
	FieldCount
	FieldSum
)

// Field database type.
type FieldDB string

// Field model.
type Field struct {
	DB     FieldDB
	Type   FieldType
	Ignore []QueryBuilderType // Ignore query builder types.
}

// NewField creates a new Field instance with the specified database field name and ignore types.
// It takes a variable number of arguments, which can be any of the following:
// FieldDB: the database field name.
// QueryBuilderType: the query type to ignore.
// The function returns the newly created Field model.
func NewField(params ...any) *Field {
	// field model
	f := &Field{Type: FieldDefault}

	// extract names
	for _, param := range params {
		switch param := param.(type) {
		case FieldDB:
			f.DB = param
		case QueryBuilderType:
			f.Ignore = append(f.Ignore, param)
		}
	}

	// return field
	return f
}

// NewAllField returns a new Field model with DB type set to "*".
//
// The returned Field model is equivalent to calling NewField("*").
func NewAllField() *Field {
	return &Field{
		DB:   "*",
		Type: FieldAll,
	}
}

// NewSumField creates new Field model with sum type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with Type set to FieldSum.
//
// Returns created Field model with sum type.
func NewSumField(field *Field) *Field {
	return &Field{
		DB:   field.DB,
		Type: FieldSum,
	}
}

// NewCountField creates a new Field model with count type.
//
// It takes the existing Field model and creates a new one with the same
// DB field and with Type set to FieldCount.
//
// Returns created Field model with count type.
func NewCountField(field *Field) *Field {
	return &Field{
		DB:   field.DB,
		Type: FieldCount,
	}
}
