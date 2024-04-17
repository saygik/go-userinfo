package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserApiTokenByName(login string) (user entity.IdName, err error) {
	sql := fmt.Sprintf(`SELECT id, ifnull(api_token,'') AS name
       FROM glpi_users  WHERE glpi_users.name= '%s'`, login)
	err = r.db.SelectOne(&user, sql)
	return user, err
}
