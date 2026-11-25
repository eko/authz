package rule

// Operator is used in Authz to identify comparison in attribute rules.
type Operator string

const (
	// OperatorEqual is the key used to identify an equality check.
	OperatorEqual Operator = "=="

	// OperatorNotEqual is the key used to identify a not equal check.
	OperatorNotEqual Operator = "!="

	// OperatorGreater is the key used to identify a numeric greater check.
	OperatorGreater Operator = ">"

	// OperatorGreaterEqual is the key used to identify a numeric greater or equal check.
	OperatorGreaterEqual Operator = ">="

	// OperatorLower is the key used to identify a numeric lower check.
	OperatorLower Operator = "<"

	// OperatorLowerEqual is the key used to identify a numeric lower or equal check.
	OperatorLowerEqual Operator = "<="

	// OperatorContains is the key used to identify a check where a comma-separated
	// list attribute value contains the given value.
	OperatorContains Operator = "contains"
)
