package rule

// Operator is used in Authz to identify comparison in attribute rules.
type Operator string

const (
	// OperatorEqual is the key used to identify an equality check.
	OperatorEqual Operator = "=="

	// OperatorNotEqual is the key used to identify a not equal check.
	OperatorNotEqual Operator = "!="
)
