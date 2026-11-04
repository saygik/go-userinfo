package usecase

import "github.com/saygik/go-userinfo/internal/entity"

// Действия аудита изменений прав.
const (
	AuditRoleAdd       = "role_add"
	AuditRoleRemove    = "role_remove"
	AuditDomainAdd     = "domain_add"
	AuditDomainUpdate  = "domain_update"
	AuditDomainRemove  = "domain_remove"
	AuditSectionAdd    = "section_add"
	AuditSectionRemove = "section_remove"
	AuditRightsClear   = "rights_clear"
)

// auditPermissionChange пишет запись аудита после успешного изменения прав.
// Ошибка аудита не отменяет уже выполненное изменение и только логируется.
func (u *UseCase) auditPermissionChange(perms entity.Permissions, target, action, object, details string) {
	if u.repo == nil {
		return
	}
	if err := u.repo.AddPermissionAudit(perms.User, target, action, object, details); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось записать аудит изменения прав")
	}
}

// GetPermissionAudit возвращает последние записи аудита изменений прав.
func (u *UseCase) GetPermissionAudit(limit int) ([]entity.PermissionAudit, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	rows, err := u.repo.GetPermissionAudits(limit)
	if rows == nil {
		rows = []entity.PermissionAudit{}
	}
	return rows, err
}
