package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserProfiles(id int) (profiles []entity.GLPIUserProfile, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_profiles_users.profiles_id AS id, glpi_profiles.name AS 'name', glpi_entities.completename AS ename, glpi_profiles_users.is_recursive AS 'recursive',
		glpi_entities.id AS 'eid', IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs'
		FROM  glpi_profiles_users
		INNER JOIN glpi_profiles ON glpi_profiles_users.profiles_id=glpi_profiles.id
		INNER JOIN glpi_entities ON glpi_profiles_users.entities_id=glpi_entities.id
		WHERE users_id=%d`, id)
	_, err = r.db.Select(&profiles, sql)
	return profiles, err
}
