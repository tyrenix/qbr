package domain

// SqlPlaceholder type.
type SqlPlaceholder string

// Sql placeholders variables.
const (
	SqlDollar   SqlPlaceholder = "$"
	SqlQuestion SqlPlaceholder = "?"
)
