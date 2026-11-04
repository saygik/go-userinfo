package usecase

import "strconv"

func (u *UseCase) DelUserSoftware(id string, editor string) error {
	recordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return u.Error("неправильный id записи пользователя системы")
	}

	record, _ := u.repo.GetUserSoftwareByID(recordID)

	if err := u.repo.DelOneUserSoftware(id); err != nil {
		return err
	}

	u.auditUserSoftwareChange(editor, UserSoftwareAuditDelete, record, userSoftwareAuditDetails(UserSoftwareAuditDelete, record))
	return nil
}
