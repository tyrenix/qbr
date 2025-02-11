package domain

// Modification type.
type ModificationType int

// Modification types.
const (
	ModificationAdd        ModificationType = iota // amount = amount + value
	ModificationSubtract                           // amount = amount - value
	ModificationMultiply                           // amount = amount * value
	ModificationDivide                             // amount = amount / value
	ModificationBitwiseAnd                         // amount = amount & value
	ModificationBitwiseOr                          // amount = amount | value
	ModificationBitwiseXor                         // amount = amount ^ value
	ModificationShiftLeft                          // amount = amount << value
	ModificationShiftRight                         // amount = amount >> value
)

// Modification model.
type Modification struct {
	Field    *Field
	Operator ModificationType
	Value    any
}
