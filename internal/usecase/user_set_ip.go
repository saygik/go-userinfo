package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) SetUserIp(userForm entity.UserActivityForm) (string, error) {

	domain := getDomainFromUserName(userForm.User)
	if !u.ad.IsDomainExist(domain) {
		return "", u.Error("такой домен отсутствует в системе")
	}

	if userForm.Activiy == "" {
		userForm.Activiy = "login"
	}

	msgResponce, err := u.repo.SetUserIp(userForm)
	if err != nil {
		return "", u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	return msgResponce, nil
}
