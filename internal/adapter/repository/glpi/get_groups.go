package glpi

import "github.com/saygik/go-userinfo/internal/entity"

// GetGroups возвращает справочник групп GLPI для выбора прав в админке календарей.
func (r *Repository) GetGroups() (groups []entity.IdName, err error) {
	_, err = r.db.Select(&groups, `
		SELECT id, name
		FROM glpi_groups
		WHERE name <> ''
		ORDER BY name`)
	return groups, err
}
