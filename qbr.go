package qbr

// Query annotation type.
type QueryAnnotationType string

// Query annotation types.
const (
	QueryQbr      QueryAnnotationType = "qbr"
	QueryDB       QueryAnnotationType = "db"
	QueryIgnoreOn QueryAnnotationType = "ignore_on"
)

// Operation type.
type OperationType string

// Operation types.
const (
	OperationCreate OperationType = "create"
	OperationRead   OperationType = "read"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
)

// Query builder model.
type QueryBuilder struct {
	selects       []Field
	conditions    []Condition
	orderBy       []OrderBy
	data          []Data
	limit         uint64
	offset        uint64
	operationType OperationType
}

// New creates new query builder with given query type.
//
// Returns created query builder.
func New(t OperationType) *QueryBuilder {
	// create and return query builder
	return &QueryBuilder{operationType: t}
}
