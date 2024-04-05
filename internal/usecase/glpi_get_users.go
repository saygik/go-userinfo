package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGLPIUsers() ([]entity.GLPIUserShort, error) {
	users, err := u.glpi.GetUsers()
	if err != nil {
		return users, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}

	for i, user := range users {
		adUser := u.GetUserADPropertysShort(user.Name)
		users[i].ADProfile = adUser
	}

	return users, nil
}
