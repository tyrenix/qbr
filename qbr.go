package qbr

import "github.com/tyrenix/qbr/domain"

// Query model.
type Query struct {
	operation  domain.OperationType
	selects    []domain.Field
	conditions []domain.Condition
	sort       []domain.Sort
	data       []domain.Data
	lock       bool
	limit      uint64
	offset     uint64
	suffix     string
}

// New creates new query builder with given query type.
//
// Returns created query builder.
func New(t domain.OperationType) *Query {
	// create and return query builder
	return &Query{
		operation: t,
		selects:   []domain.Field{*NewAllField()},
	}
}

// NewCreate creates a new query builder with OperationCreate type.
//
// Returns the created query builder.
func NewCreate() *Query {
	return New(domain.OperationCreate)
}

// NewRead creates a new query builder with OperationRead type.
//
// Returns the created query builder.
func NewRead() *Query {
	return New(domain.OperationRead)
}

// NewUpdate creates new query builder with OperationUpdate type.
//
// Returns created query builder.
func NewUpdate() *Query {
	return New(domain.OperationUpdate)
}

// NewDelete creates new query builder with OperationDelete type.
//
// Returns created query builder.
func NewDelete() *Query {
	return New(domain.OperationDelete)
}
