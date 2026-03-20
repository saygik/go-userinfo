package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetSoftwaresAdmins() (admins []entity.SoftwareAdmins, err error) {
	sql := `
	SELECT DISTINCT
    gu.groups_id AS id,
    u.name AS name
FROM glpi_groups_users gu
LEFT JOIN glpi_users u ON u.id = gu.users_id
WHERE gu.groups_id IN (
    SELECT DISTINCT gi.groups_id
    FROM glpi_groups_items gi
    WHERE gi.itemtype = 'Software'
      AND gi.type = 2
      AND gi.items_id IN (SELECT id FROM glpi_softwares)
)
ORDER BY id;
	`
	_, err = r.db.Select(&admins, sql)
	return admins, err
}

func (r *Repository) GetSoftwareAdmins(id int) (admins []entity.SoftwareAdmins, err error) {
	sql := fmt.Sprintf(`
	SELECT DISTINCT
    gu.groups_id AS id,
    u.name AS name
FROM glpi_groups_users gu
LEFT JOIN glpi_users u ON u.id = gu.users_id
WHERE gu.groups_id IN (
    SELECT DISTINCT gi.groups_id
    FROM glpi_groups_items gi
    WHERE gi.itemtype = 'Software'
      AND gi.type = 2
      AND gi.items_id = %d
)
ORDER BY id;
	`, id)
	_, err = r.db.Select(&admins, sql)
	return admins, err
}
