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
	domain.OperatorIn:                 "IN",
}

// sqlModifications is a map that defines SQL modifications for different ModificationTypes.
// It currently supports all supported modifications.
var sqlModifications = map[domain.ModificationType]string{
	domain.ModificationAdd:        "+",
	domain.ModificationSubtract:   "-",
	domain.ModificationMultiply:   "*",
	domain.ModificationDivide:     "/",
	domain.ModificationBitwiseAnd: "&",
	domain.ModificationBitwiseOr:  "|",
	domain.ModificationBitwiseXor: "^",
	domain.ModificationShiftLeft:  "<<",
	domain.ModificationShiftRight: ">>",
}
