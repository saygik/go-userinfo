package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetAllSchedules(currentUser string) ([]entity.IdName, error) {
	glpiGroups := u.glpi.GetUserGroupsListByName(currentUser)

	schedules, err := u.repo.GetScheduleUserAvailableGroups(currentUser, glpiGroups)
	if err != nil {
		return nil, u.Error("невозможно получить календари пользователя из базы данных")
	}

	return schedules, nil
}
