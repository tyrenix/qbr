package qbr

// Limit set limit.
func (qb *Query) Limit(limit uint64) *Query {
	// set limit
	qb.limit = limit

	// return query
	return qb
}

// GetLimit returns the limit set for the query, or 0 if no limit has been set.
func (qb *Query) GetLimit() uint64 {
	return qb.limit
}
