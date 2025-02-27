package qbr

// Suffix adds a suffix to the query builder.
func (q *Query) Suffix(s string) *Query {
	q.suffix = s
	return q
}

// GetSuffix returns the suffix of the query builder.
func (q *Query) GetSuffix() string {
	return q.suffix
}
