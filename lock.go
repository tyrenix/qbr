package qbr

// Lock sets the FOR UPDATE lock on the query, which causes the rows returned by
// SELECT statement to be locked as though for update. This means that the rows
// are locked until the end of the transaction, regardless of the setting of the
// READ COMMITTED isolation level.
//
// Use this method to set the lock on the query. Returns the modified QueryBuilder
// instance for method chaining.
func (qb *Query) Lock() *Query {
	qb.lock = true
	return qb
}

// IsLock returns true if the query has been set with a FOR UPDATE lock,
// indicating that the rows returned by the SELECT statement are locked for
// update until the end of the transaction. Otherwise, it returns false.
func (qb *Query) IsLock() bool {
	return qb.lock
}
