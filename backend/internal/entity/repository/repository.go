package repository

import "github.com/eko/authz/backend/internal/entity/model"

// Repositories is a constraint interface that allows only authz library repositories.
type Repositories interface {
	base[model.Action] |
		base[model.Client] |
		base[model.CompiledPolicy] |
		base[model.Policy] |
		base[model.Principal] |
		base[model.Resource] |
		base[model.Role] |
		base[model.User] |
		resource
}
