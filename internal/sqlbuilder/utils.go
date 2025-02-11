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

// valueToDBValue takes a value and returns a value that can be used in a
// SQL query. If the value is a ValueType, it returns an error if it is not a
// null value. If the value is a struct or a pointer to a struct, it marshals
// the struct to JSON and returns it as a string. Otherwise, it returns the
// original value.
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

// buildSetData formats a Data object into a SQL SET data string, along with a placeholder, value, and error if the value could not be converted.
// It takes a Data object, a domain.SqlPlaceholder for parameter substitution, and a parameter index.
// If the Data object's Value is a Modification, it extracts the Modification, gets the corresponding SQL operator, and converts the Modification's Value to a database-compatible value.
// If the Data object's Value is not a Modification, it simply converts the value to a database-compatible value.
// The function returns the database field name, the placeholder string, the converted value, and an error if the value could not be converted.
func buildSetData(data domain.Data, placeholder domain.SqlPlaceholder, index int) (field, plc string, value any, err error) {
	// extract value
	value = data.Value

	// is value is modification
	if mod, ok := value.(*domain.Modification); ok {
		// get operator
		op, exists := sqlModifications[mod.Operator]
		if !exists {
			return "", "", nil, fmt.Errorf("unsupported modification operator: %d", mod.Operator)
		}

		// convert internal value
		inner, err := valueToDBValue(mod.Value)
		if err != nil {
			return "", "", nil, err
		}

		// create placeholder
		plc = fmt.Sprintf("%s %s %v", getFieldName(mod.Field), op, getPlaceholder(placeholder, index))

		// return data and create placeholder
		return getFieldName(data.Field), plc, inner, nil
	}

	// create field
	field = getFieldName(data.Field)

	// create placeholder
	plc = getPlaceholder(placeholder, index)

	// convert value
	value, err = valueToDBValue(value)
	if err != nil {
		return "", "", nil, err
	}

	// return data and create placeholder
	return field, plc, value, nil
}
