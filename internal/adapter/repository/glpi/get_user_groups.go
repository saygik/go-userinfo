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

func (r *Repository) GetUserGroupsListByName(user string) (list string) {
	sql := fmt.Sprintf(
		`SELECT GROUP_CONCAT(glpi_groups.id SEPARATOR ',') AS ids  from glpi_groups
		INNER JOIN glpi_groups_users ON glpi_groups_users.groups_id=glpi_groups.id
		INNER Join glpi_users ON glpi_groups_users.users_id=glpi_users.id
		WHERE glpi_users.name='%s'`, user)
	err := r.db.SelectOne(&list, sql)
	if err != nil {
		return "0"
	}
	return list
}
