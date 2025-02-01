package qbr

import "github.com/tyrenix/qbr/domain"

// GetOperation returns the operation type of the query builder.
func (qb *Query) GetOperation() domain.OperationType {
	return qb.operation
}
