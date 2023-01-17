package manager

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
)

type AuditRepository repository.Base[model.Audit]

type Audit interface {
	BatchAdd(audits []*model.Audit) error
	GetRepository() AuditRepository
}

type auditManager struct {
	repository AuditRepository
}

// NewAudit initializes a new audit manager.
func NewAudit(
	repository AuditRepository,
) Audit {
	return &auditManager{
		repository: repository,
	}
}

func (m *auditManager) GetRepository() AuditRepository {
	return m.repository
}

func (m *auditManager) BatchAdd(audits []*model.Audit) error {
	return m.repository.Create(audits...)
}
