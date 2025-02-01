package domain

// Sort sort type.
type SortType string

// Sort types.
const (
	SortAsc  SortType = "ASC"
	SortDesc SortType = "DESC"
)

// Sort model
type Sort struct {
	Field *Field
	Type  SortType
}
