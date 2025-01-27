package qbr

// Offset set offset.
func (qb *QueryBuilder) Offset(offset uint64) *QueryBuilder {
	// set offset
	qb.offset = offset

	// return query builder
	return qb
}

// GetOffset returns the offset set for the query, or 0 if no offset has been set.
func (qb *QueryBuilder) GetOffset() uint64 {
	return qb.offset
}
