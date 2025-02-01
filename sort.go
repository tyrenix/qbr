package qbr

import "github.com/tyrenix/qbr/domain"

// NewSortDesc creates a new instance of domain.Sort with descending sort type.
//
// It takes a pointer to a domain.Field as an argument and returns a pointer
// to the newly created domain.Sort object with the provided field and
// descending sort type.
func NewSortDesc(field *domain.Field) *domain.Sort {
	return &domain.Sort{
		Field: field,
		Type:  domain.SortDesc,
	}
}

// NewSortAsc creates a new instance of domain.Sort with ascending sort type.
//
// It takes a pointer to a domain.Field as an argument and returns a pointer
// to a domain.Sort object with the specified field and SortType set to SortAsc.
func NewSortAsc(field *domain.Field) *domain.Sort {
	return &domain.Sort{
		Field: field,
		Type:  domain.SortAsc,
	}
}

// Sort add sort.
func (qb *Query) Sort(sorts ...*domain.Sort) *Query {
	// add sorts to query
	for _, sort := range sorts {
		qb.sort = append(qb.sort, *sort)
	}

	// return query
	return qb
}

// GetSort returns the sort parameters of the query, or an empty slice if no order by has been set.
func (qb *Query) GetSort() []domain.Sort {
	return qb.sort
}
