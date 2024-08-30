package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetUserSoftwaresByFio(fio string) ([]entity.Software, error) {

	if fio == "" {
		return nil, u.Error("ФИО пользователя в запросе отсутствует")
	}

	softwares, err := u.glpi.GetAllSoftwares()
	if err != nil {
		return nil, u.Error("невозможно получить список систем из GLPI")
	}

	userSoftwares, err := u.repo.GetUserSoftwaresByFio(fio)
	if err != nil {
		return nil, u.Error("невозможно получить список систем пользователя")
	}
	filteredSoft := []entity.Software{}
	for _, soft := range softwares {
		for _, idsoft := range userSoftwares {
			if soft.Id == idsoft {
				filteredSoft = append(filteredSoft, soft)
			}
		}
	}

	return filteredSoft, nil

}
