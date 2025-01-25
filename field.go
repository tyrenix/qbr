package qbr

// Field type.
type FieldType int

const (
	FieldDefault FieldType = iota
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

// NewField creates new Field model from given names.
//
// NewField accepts variable arguments of FieldDBType and FieldJSONType.
// It iterates over the arguments and sets DB and JSON fields of the model.
// If the argument is not recognized, it skips it.
//
// Returns created Field model.
func NewField(names ...any) *Field {
	// field model
	f := &Field{Type: FieldDefault}

	// extract names
	for _, name := range names {
		switch name := name.(type) {
		case FieldDB:
			f.DB = name
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
		Type: FieldDefault,
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
