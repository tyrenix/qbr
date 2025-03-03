package qbr

import "github.com/tyrenix/qbr/domain"

// FilterMatchingConditions recursively filters and returns all conditions
// that have a field equal to the targetField. It traverses nested conditions.
func FilterMatchingConditions(targetField *domain.Field, conds []domain.Condition) []domain.Condition {
	var result []domain.Condition
	for _, cond := range conds {
		// if condition value is a slice of nested conditions, process them recursively.
		if nested, ok := cond.Value.([]domain.Condition); ok {
			result = append(result, FilterMatchingConditions(targetField, nested)...)
			continue
		}
		// if condition's field matches the target field, add to result.
		if IsFieldEqual(targetField, cond.Field) {
			result = append(result, cond)
		}
	}
	return result
}
