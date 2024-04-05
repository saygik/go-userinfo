package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSoftwares() ([]entity.Software, error) {

	softwares, err := u.glpi.GetAllSoftwares()
	if err != nil {

		return nil, u.Error("невозможно получить список систем из GLPI")
	}

	// список администраторов систем
	admins, _ := u.glpi.GetSoftwaresAdmins()
	if err != nil {
		return nil, u.Error("невозможно получить список администраторов систем из GLPI")
	}
	softAdmins := []map[string]interface{}{}
	for i, soft := range softwares {
		for _, admin := range admins {
			if soft.Groups_id_tech == admin.Id {
				adUser := u.GetUserADPropertysShort(admin.Name)
				softAdmins = append(softAdmins, adUser)
			}
		}
		if len(softAdmins) > 0 {
			softwares[i].Admins = softAdmins
			softAdmins = []map[string]interface{}{}
		} else {
			soft.Admins = []map[string]interface{}{}
		}

	}

	return softwares, nil

}
