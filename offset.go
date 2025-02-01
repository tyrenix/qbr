package qbr

// Offset set offset.
func (qb *Query) Offset(offset uint64) *Query {
	// set offset
	qb.offset = offset

	// return query
	return qb
}

// GetOffset returns the offset set for the query, or 0 if no offset has been set.
func (qb *Query) GetOffset() uint64 {
	return qb.offset
}
