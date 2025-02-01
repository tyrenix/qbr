package domain

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
