package model

// Models is a constraint interface that allows only authz library models.
type Models interface {
	Action | Attribute | CompiledPolicy | Policy | Resource | Role | Principal
}
