package usecase

import (
	"fmt"
)

func (u *UseCase) DelScheduleTask(id string) error {
	err := u.repo.DelScheduleTask(id)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	return nil
}
