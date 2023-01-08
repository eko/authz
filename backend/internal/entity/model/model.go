package model

// Models is a constraint interface that allows only authz library models.
type Models interface {
	Action | Attribute | Client | CompiledPolicy | Policy | Principal | Resource | Role | Token | User
}
