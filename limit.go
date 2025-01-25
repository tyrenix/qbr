package qbr

// Limit set limit.
func (qb *QueryBuilder) Limit(limit uint64) *QueryBuilder {
	// set limit
	qb.limit = limit

	// return query builder
	return qb
}
