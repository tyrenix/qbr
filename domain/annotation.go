package domain

// Query annotation type.
type QueryAnnotationType string

// Query annotation types.
const (
	QueryQbr      QueryAnnotationType = "qbr"
	QueryDB       QueryAnnotationType = "db"
	QueryIgnoreOn QueryAnnotationType = "ignore_on"
)
