package sqlbuilder

import "github.com/tyrenix/qbr/domain"

// sqlAggregationFormats is a map that defines SQL aggregation formats for different AggregationTypes.
// It currently supports all supported aggregation types.
var sqlAggregationFormats = map[domain.AggregationType]string{
	domain.AggregationNone:  "%s",
	domain.AggregationCount: "COUNT(%s)",
	domain.AggregationSum:   "SUM(%s)",
}

// sqlOperators is a map that defines SQL operators for different OperatorTypes.
// It currently supports all supported operators.
var sqlOperators = map[domain.OperatorType]string{
	domain.OperatorAnd:                "AND",
	domain.OperatorOr:                 "OR",
	domain.OperatorEqual:              "=",
	domain.OperatorNotEqual:           "!=",
	domain.OperatorLessThan:           "<",
	domain.OperatorGreaterThan:        ">",
	domain.OperatorLessThanOrEqual:    "<=",
	domain.OperatorGreaterThanOrEqual: ">=",
}
