package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetUserADActivity(userName string, userTechName string) ([]entity.UserActivity, error) {

	if userName == "" {
		return nil, u.Error("имя пользователя в запросе отсутствует")
	}
	if !isEmailValid(userName) {
		return nil, u.Error("неверное имя пользователя в запросе")
	}

	if userTechName == "" {
		return nil, u.Error("сначала войдите в систему")
	}

	domain := getDomainFromUserName(userName)
	access := u.GetAccessToResource(domain, userTechName)
	//	accessToTechnicalInfo := access == 1
	if !(access == 1) {
		return nil, u.Error("у вас нет прав на просмотр этой информации в домене")
	}

	activities, err := u.repo.GetUserActivity(userName)
	if err != nil {

		return nil, u.Error("невозможно получить информацию об активности пользователя с сервера")
	}

	return activities, nil

}
