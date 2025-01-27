package qbr

// Limit set limit.
func (qb *QueryBuilder) Limit(limit uint64) *QueryBuilder {
	// set limit
	qb.limit = limit

	// return query builder
	return qb
}

// GetLimit returns the limit set for the query, or 0 if no limit has been set.
func (qb *QueryBuilder) GetLimit() uint64 {
	return qb.limit
}
