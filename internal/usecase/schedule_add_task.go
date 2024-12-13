package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddScheduleTask(task entity.ScheduleTask) (entity.ScheduleTaskCalendar, error) {
	task, err := u.repo.AddScheduleTask(task)
	if err != nil {
		return entity.ScheduleTaskCalendar{}, u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	oneTask := u.scheduleTaskToScheduleTaskCalendar(task)
	return oneTask, nil
}
