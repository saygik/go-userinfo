package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetScheduleTasks(id string) ([]entity.ScheduleTask, error) {

	schedule, err := u.repo.GetScheduleTasks(id)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	return schedule, nil
}
