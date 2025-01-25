package qbr

// Query builder type.
type QueryBuilderType int

// Query builder types.
const (
	QueryBuilderCreate QueryBuilderType = iota
	QueryBuilderRead
	QueryBuilderUpdate
	QueryBuilderDelete
)

// Query builder model.
type QueryBuilder struct {
	selects    []Field
	conditions []Condition
	orderBy    []OrderBy
	data       []Data
	limit      uint64
	offset     uint64
	queryType  QueryBuilderType
}

// New creates new query builder with given query type.
//
// Returns created query builder.
func New(t QueryBuilderType) *QueryBuilder {
	// create and return query builder
	return &QueryBuilder{queryType: t}
}
