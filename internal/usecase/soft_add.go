package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddUserSoftware(softwareForm entity.SoftwareForm, editor string) error {
	if err := u.repo.AddOneUserSoftware(softwareForm); err != nil {
		return err
	}
	record := entity.SoftUser{
		IdSoft: int(softwareForm.Id),
		Name:   softwareForm.User,
	}
	u.auditUserSoftwareChange(editor, UserSoftwareAuditCreate, record, userSoftwareAuditDetails(UserSoftwareAuditCreate, record))
	return nil
}
