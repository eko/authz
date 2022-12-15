package model

// Models is a constraint interface that allows only authz library models.
type Models interface {
	Action | CompiledPolicy | Policy | Resource | Role | Principal
}
