package qbr

// Query builder annotation.
type QueryBuilderAnnotationType string

const (
	QueryBuilderIgnore QueryBuilderAnnotationType = "qbrIgnore"
	QueryBuilderDB     QueryBuilderAnnotationType = "db"
)

// Query builder type.
type QueryBuilderType string

// Query builder types.
const (
	QueryBuilderCreate QueryBuilderType = "create"
	QueryBuilderRead   QueryBuilderType = "read"
	QueryBuilderUpdate QueryBuilderType = "update"
	QueryBuilderDelete QueryBuilderType = "delete"
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
