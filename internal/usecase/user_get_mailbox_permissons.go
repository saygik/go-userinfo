package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetUserMailboxPermissions(userName string, userTechName string) ([]entity.MailBoxDelegates, error) {

	if userName == "" {
		return nil, u.Error("имя пользователя в запросе отсутствует")
	}
	if !isEmailValid(userName) {
		return nil, u.Error("неверное имя пользователя в запросе")
	}

	if userTechName == "" {
		return nil, u.Error("сначала войдите в систему")
	}
	if !u.HasTechnicalRole(userTechName) {
		return nil, u.Error("у вас нет прав на получение прав пользователя на почтовые ящики")
	}
	return u.repo.GetOneDelegateMailBoxes(userName)

}
