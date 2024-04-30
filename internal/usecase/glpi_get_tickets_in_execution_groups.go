package usecase

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetTicketsInExecutionGroups(user string) ([]entity.GLPI_Ticket, error) {
	groups, err := u.repo.GetUserGlpiTrackingGroups(user)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}

	str := strings.Replace(strings.Replace(strings.Replace(strings.Trim(fmt.Sprint(groups), "[]"), " ", ",", -1), "{", "", -1), "}", "", -1)
	if len(str) < 1 {
		return nil, u.Error("нет групп GLPI для слежения для этого пользователя")
	}

	tickets, err := u.glpi.GetTicketsInExecutionGroups(str)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return tickets, nil
}
