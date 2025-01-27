package qbr

// Order by sort type.
type OrderByType string

// Order by types.
const (
	OrderByAsc  OrderByType = "ASC"
	OrderByDesc OrderByType = "DESC"
)

// Order by model
type OrderBy struct {
	Field *Field
	Type  OrderByType
}

// OrderBy add sort.
func (qb *QueryBuilder) OrderBy(params ...*OrderBy) *QueryBuilder {
	// add sort params to query
	for _, param := range params {
		qb.orderBy = append(qb.orderBy, *param)
	}

	// return query builder
	return qb
}

// GetOrderBy returns the order by parameters of the query, or an empty slice if no order by has been set.
func (qb *QueryBuilder) GetOrderBy() []OrderBy {
	return qb.orderBy
}
