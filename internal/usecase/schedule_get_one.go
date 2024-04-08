package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSchedule(id string) (entity.Schedule, error) {

	schedule, err := u.repo.GetSchedule(id)
	if err != nil {
		return entity.Schedule{}, u.Error("имя пользователя в запросе отсутствует")
	}
	return schedule, nil
}
