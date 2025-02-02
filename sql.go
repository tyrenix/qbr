package qbr

import (
	"fmt"

	"github.com/tyrenix/qbr/domain"
	"github.com/tyrenix/qbr/internal/sqlbuilder"
)

// SqlPlaceholder is a placeholder type for SQL query placeholders.
const (
	SqlDollar   domain.SqlPlaceholder = "$"
	SqlQuestion domain.SqlPlaceholder = "?"
)

// ToSql builds SQL query from the query builder data and returns it as a string, along with the query parameters and an error if the query could not be built.
//
// It supports the following query types: SELECT, INSERT, UPDATE, DELETE.
func (qb *Query) ToSql(table string, placeholder domain.SqlPlaceholder) (string, []any, error) {
	// select need method for build
	switch qb.operation {
	case domain.OperationRead:
		return sqlbuilder.CreateSelectSql(qb, table, placeholder)
	case domain.OperationCreate:
		return sqlbuilder.CreateInsertSql(qb, table, placeholder)
	case domain.OperationUpdate:
		return sqlbuilder.CreateUpdateSql(qb, table, placeholder)
	case domain.OperationDelete:
		return sqlbuilder.CreateDeleteSql(qb, table, placeholder)
	default:
		return "", nil, fmt.Errorf("unsupported query type: %v", qb.operation)
	}
}
