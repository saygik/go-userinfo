package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserGroups(id int) (groups []entity.IdName, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_groups_users.groups_id AS id, glpi_groups.name AS name
		from glpi_groups_users INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_users.groups_id
		WHERE users_id=%d`, id)
	_, err = r.db.Select(&groups, sql)
	return groups, err
}
