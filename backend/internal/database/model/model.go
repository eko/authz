package model

// Models is a constraint interface that allows only authz library models.
type Models interface {
	Action | Policy | Resource | Role | Subject
}
