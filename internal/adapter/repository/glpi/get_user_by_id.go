package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserById(id int) (user entity.GLPIUser, err error) {
	sql := fmt.Sprintf(
		`SELECT name, CONCAT(realname," ", firstname) AS fio, IFNULL((SELECT email from glpi_useremails WHERE users_id=glpi_users.id LIMIT 1),'') as email
         FROM glpi_users WHERE id=%d`, id)
	err = r.db.SelectOne(&user, sql)
	return user, err
}
