package qbr

import (
	"reflect"
	"strings"
	"time"

	"github.com/tyrenix/qbr/domain"
)

// isZero check on zero value
func isZero(value any) bool {
	// check is nil
	if value == nil {
		return true
	}

	// get value by reflect
	v := reflect.ValueOf(value)

	// for pointer types
	if v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Interface {
		return v.IsNil()
	}

	// is zero
	if v.IsZero() {
		return true
	}

	// for Time types
	if v.Kind() == reflect.Struct && v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	// for other types return false
	return false
}

// isFieldIgnored checks if a field is ignored for a given query type.
//
// The function checks if the query type is in the field's list of ignored operations.
// If it is, the function returns true, indicating that the field is ignored. Otherwise,
// it returns false.
func isFieldIgnored(field *domain.Field, queryType domain.OperationType) bool {
	// check is ignored
	for _, ignoreOp := range field.IgnoreOn {
		if ignoreOp == queryType {
			return true
		}
	}

	// not ignored
	return false
}

// extractFieldFromStruct extracts a Field object from a given struct field.
//
// The function retrieves the "db" tag from the field annotation and uses it to
// initialize a Field object. If the "db" tag is empty, the function returns nil.
// Additionally, the function checks for a "qbr" tag and parses any annotations
// it contains. If the "qbr" tag includes an "ignore_on" annotation, the function
// extracts the ignored operations and adds them to the Field's IgnoredOperations
// slice.
//
// The resulting Field object is returned, representing a database field with
// optional ignored operations based on the struct field's annotations.
func extractFieldFromStruct(ft reflect.StructField) *domain.Field {
	// get tags from field annotation
	db := ft.Tag.Get(string(domain.QueryDB))

	// check is not empty
	if db == "" {
		return nil
	}

	// create field
	field := &domain.Field{
		DB: db,
	}

	// query builder tag
	qbr := ft.Tag.Get(string(domain.QueryQbr))

	// check is not empty
	if qbr == "" {
		return field
	}

	// get annotations from query builder annotation
	for _, block := range strings.Split(qbr, " ") {
		// check is not empty
		if block == "" {
			continue
		}

		// get annotation
		switch {
		case strings.HasPrefix(block, string(domain.QueryIgnoreOn)+"="):
			field.IgnoreOn = append(
				field.IgnoreOn,
				extractIgnoredOperationOnAnnotations(block)...,
			)
		default:
			continue
		}
	}

	// return fields
	return field
}

// extractDataFromStruct extracts fields from a given struct and returns them as a slice of Data.
// If the input is a pointer, it dereferences it before processing. The function checks if the input
// is a valid struct type and iterates through its fields. For each field, it retrieves the field's
// value and annotation, and constructs a Data object. Fields with a nil value or that do not have
// a "db" annotation are ignored. The resulting slice of Data objects is returned, representing the
// struct's fields ready for inclusion in a query.
func extractDataFromStruct(s any) []*domain.Data {
	// struct value
	val := reflect.ValueOf(s)
	// struct type
	t := reflect.TypeOf(s)

	// check is pointer
	if val.Kind() == reflect.Ptr {
		// check is nil
		if val.IsNil() {
			return nil
		}

		// dereference pointer
		val = val.Elem()
		t = t.Elem()
	}

	// check is struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	// create data slice
	var data []*domain.Data

	// we go through the fields of the structure
	for i := 0; i < val.NumField(); i++ {
		// struct field
		field := val.Field(i)
		// field type
		ft := t.Field(i)

		// add data
		data = append(data, NewData(
			extractFieldFromStruct(ft),
			field.Interface(),
			false,
		))
	}

	// return query
	return data
}

// extractIgnoredOperationOnAnnotations extracts the ignored operations from the given block string.
//
// The block string is expected to be in the format "ignore_on=<operation1>,<operation2>,...".
//
// The function splits the block by comma, trims the resulting strings, and adds them to a slice of
// ignored operations. The operation types are converted to lower case to ensure consistency.
//
// The function returns the slice of ignored operations.
func extractIgnoredOperationOnAnnotations(block string) []domain.OperationType {
	// delete from block annotation type
	block = strings.TrimPrefix(block, string(domain.QueryIgnoreOn)+"=")

	// split by comma
	ops := strings.Split(block, ",")

	// slice of ignored operations
	ignOps := make([]domain.OperationType, 0, len(ops))

	// add ignored operations
	for _, op := range ops {
		// check is not empty
		if op == "" {
			continue
		}

		// get operation type
		ignOps = append(ignOps, domain.OperationType(strings.ToLower(op)))
	}

	// return ignored operations
	return ignOps
}

// removeZeroCondition takes a variable number of conditions and returns a new slice
// with the following changes:
//  1. Conditions with a Value of nil or a zero value are removed.
//  2. Conditions with a Field that is ignored for the current query type are removed.
//  3. Conditions with a Value of domain.ValueNull are removed if the condition is not
//     an aggregation or an equality/inequality check.
//
// The method returns the modified slice of conditions.
func removeZeroCondition(conds ...domain.Condition) []domain.Condition {
	// conditions for return
	result := []domain.Condition{}

	// stack
	stack := append([]domain.Condition{}, conds...)

	// check all conditions
	for len(stack) > 0 {
		// condition
		cond := stack[0]
		// remove condition from stack
		stack = stack[1:]

		// select condition type
		switch v := cond.Value.(type) {
		case []domain.Condition:
			// set new removed conditions
			cond.Value = removeZeroCondition(v...)

			// add formatted conditions
			result = append(result, cond)
		default:
			// check is not ignored
			if isFieldIgnored(cond.Field, domain.OperationRead) {
				continue
			}

			// check if system conditional and handle value null case
			switch t := cond.Value.(type) {
			case domain.ValueType:
				if t == domain.ValueNull {
					// skip if value is null and not supported aggregation or operator
					if cond.Field.Aggregation != domain.AggregationNone ||
						(cond.Operator != domain.OperatorEqual && cond.Operator != domain.OperatorNotEqual) {
						continue
					}
				}
			}

			// // check is not zero
			// if !isZero(v) {
			// 	// add condition
			// 	result = append(result, cond)
			// }

			// is nil, set null value
			if v == nil {
				cond.Value = domain.ValueNull
			}

			// add condition
			result = append(result, cond)
		}
	}

	// return conditions
	return result
}
