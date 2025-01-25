package qbr

// Offset set offset.
func (qb *QueryBuilder) Offset(offset uint64) *QueryBuilder {
	// set offset
	qb.offset = offset

	// return query builder
	return qb
}
