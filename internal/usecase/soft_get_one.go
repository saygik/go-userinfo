package usecase

import (
	"strconv"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSoftware(ids string) (entity.Software, error) {
	software := entity.Software{}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return software, u.Error("неправильный id системы")
	}

	software, err = u.glpi.GetSoftware(id)
	if err != nil {
		return software, u.Error("невозможно получить систему из GLPI")
	}

	// список администраторов систем
	admins, _ := u.glpi.GetSoftwaresAdmins()
	if err != nil {
		return software, u.Error("невозможно получить список администраторов систем из GLPI")
	}
	softAdmins := []map[string]interface{}{}

	for _, admin := range admins {
		if software.Groups_id_tech == admin.Id {
			adUser := u.GetUserADPropertysShort(admin.Name)
			softAdmins = append(softAdmins, adUser)
		}
	}
	if len(softAdmins) > 0 {
		software.Admins = softAdmins
		softAdmins = []map[string]interface{}{}
	} else {
		software.Admins = []map[string]interface{}{}
	}

	return software, nil

}
