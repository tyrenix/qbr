package qbr

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

// Sort add sort.
func (qb *QueryBuilder) Sort(params ...*Sort) *QueryBuilder {
	// add sort params to query
	for _, param := range params {
		qb.sort = append(qb.sort, *param)
	}

	// return query builder
	return qb
}

// GetSort returns the sort parameters of the query, or an empty slice if no order by has been set.
func (qb *QueryBuilder) GetSort() []Sort {
	return qb.sort
}

// NewSort creates a new Sort model.
//
// It takes a Field model and a SortType (either SortAsc or SortDesc) as parameters, and returns a new Sort model
// with the given Field and SortType.
//
// The returned Sort model is ready to use and can be passed to the Sort method of a QueryBuilder in order to add
// sorting to a query.
func NewSort(field *Field, sortType SortType) *Sort {
	return &Sort{
		Field: field,
		Type:  sortType,
	}
}
