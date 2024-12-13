package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetScheduleTasks(id int) ([]entity.ScheduleTaskCalendar, error) {

	dbTasks, err := u.repo.GetScheduleTasks(id)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	tasks := []entity.ScheduleTaskCalendar{}
	for _, task := range dbTasks {
		oneTask := u.scheduleTaskToScheduleTaskCalendar(task)
		tasks = append(tasks, oneTask)
	}
	return tasks, nil
}
