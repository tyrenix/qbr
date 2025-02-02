package sqlbuilder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tyrenix/qbr/domain"
)

// getSqlOperator returns the SQL operator associated with the given OperatorType.
// If the operator is not supported, it returns an empty string.
func getSqlOperator(op domain.OperatorType) string {
	// get operator
	v, ok := sqlOperators[op]
	if !ok {
		return ""
	}

	// return operator
	return v
}

// getFieldName takes a Field object and returns the string value of its DB
// field. This is the field name in the database that the field corresponds to.
func getFieldName(field *domain.Field) string {
	return string(field.DB)
}

// getPlaceholder generates a SQL placeholder string based on the specified
// placeholder type and index. If the placeholder type is ToSqlDollar, it returns
// a parameterized string using the dollar sign format (e.g., $1, $2). Otherwise,
// it returns the placeholder type as a string.
func getPlaceholder(plc domain.SqlPlaceholder, index int) string {
	// dollar placeholder
	if plc == domain.SqlDollar {
		return fmt.Sprintf("$%d", index)
	}

	// return default placeholder
	return string(plc)
}

// valueToDBValue takes a value and returns a value that can be used in a SQL query
// and an error if the value is not supported.
//
// If the value is of type ValueType, it returns the string "NULL" if the value
// is ValueNull, otherwise it returns an error.
//
// If the value is a struct or a pointer to a struct, it converts the value to
// a JSON string and returns it.
//
// If the value is not a struct or a pointer to a struct, it returns the original
// value.
func valueToDBValue(value any) (any, error) {
	// is value is ValueType
	if v, ok := value.(domain.ValueType); ok {
		// is null value return null
		if v == domain.ValueNull {
			return nil, nil
		}

		// return nil
		return nil, fmt.Errorf("unsupported value type: %d", v)
	}

	// get reflect value
	v := reflect.ValueOf(value)

	// check if the value is a struct or a pointer to a struct
	if v.Kind() == reflect.Ptr &&
		v.Elem().Kind() == reflect.Struct {
		// convert struct to JSON
		j, err := json.Marshal(v.Elem().Interface())
		if err != nil {
			return nil, err
		}

		// return json string
		return j, nil
	}

	// return the original value if not a struct or pointer to a struct
	return value, nil
}

// buildSelectFields formats a slice of Field objects into a SQL select statement string.
// It iterates over the provided fields, and for each field, it checks if there is an
// associated SQL format in the sqlFieldFormats map based on the field's type. If a format
// exists, it retrieves the database field name and applies the format, adding the result
// to the list of select fields. The function returns a comma-separated string of the
// formatted select fields.
func buildSelects(fields []domain.Field) string {
	// fields
	var result []string

	// iterate over the slice of fields
	for _, field := range fields {
		// check is contains in map
		format, ok := sqlAggregationFormats[field.Aggregation]
		if !ok {
			continue
		}

		// append the formatted field to the result slice
		result = append(
			result,
			fmt.Sprintf(
				format,
				getFieldName(&field), // get database field name
			), // create field name
		)
	}

	// return the fields as a comma-separated string
	return strings.Join(result, ", ")
}
