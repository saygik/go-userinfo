package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) UpdateScheduleTask(task entity.ScheduleTask) error {
	err := u.repo.UpdateScheduleTask(task)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	return nil
}
