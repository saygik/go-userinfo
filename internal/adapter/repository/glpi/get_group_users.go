package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetGroupUsers(id int) (users []entity.IdName, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_groups_users.users_id AS id, glpi_users.name AS name
		from glpi_groups_users INNER JOIN glpi_users ON glpi_users.id=glpi_groups_users.users_id
		WHERE glpi_users.name NOT LIKE '%%@local%%' and glpi_groups_users.groups_id=%d`, id)
	_, err = r.db.Select(&users, sql)
	return users, err
}
