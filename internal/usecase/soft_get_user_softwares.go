package usecase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/saygik/go-userinfo/models"
)

func (u *UseCase) GetUserSoftwares(userName string) ([]entity.Software, error) {

	if userName == "" {
		return nil, u.Error("имя пользователя в запросе отсутствует")
	}
	if !isEmailValid(userName) {
		return nil, u.Error("неверное имя пользователя в запросе")
	}

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

	userSoftwares, err := userIPModel.GetUserSoftwares(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список систем пользователя"})
		return
	}
	filteredSoft := []models.Software{}
	for _, soft := range softwares {
		for _, idsoft := range userSoftwares {
			if soft.Id == idsoft {
				filteredSoft = append(filteredSoft, soft)
			}
		}
	}

	return filteredSoft, nil

}
