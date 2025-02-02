package sqlbuilder

import "github.com/tyrenix/qbr/domain"

type Query interface {
	GetSelects() []domain.Field
	GetConditions() []domain.Condition
	GetData() []domain.Data
	GetSort() []domain.Sort
	GetLimit() uint64
	GetOffset() uint64
}
