package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetSoftwaresAdmins() (admins []entity.SoftwareAdmins, err error) {
	sql := `SELECT  glpi_groups_users.groups_id AS 'id', glpi_users.name AS 'name' FROM glpi_groups_users
		 LEFT JOIN glpi_users ON glpi_users.id=glpi_groups_users.users_id
		 WHERE glpi_groups_users.groups_id IN (SELECT DISTINCT glpi_softwares.groups_id_tech FROM glpi_softwares)`
	_, err = r.db.Select(&admins, sql)
	return admins, err
}
