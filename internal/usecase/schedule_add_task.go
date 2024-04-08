package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddScheduleTask(task entity.ScheduleTask) (entity.ScheduleTask, error) {
	task, err := u.repo.AddScheduleTask(task)
	if err != nil {
		return task, u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	return task, nil
}
