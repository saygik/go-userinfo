package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetSoftware(id int) (software entity.Software, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_softwares.id, glpi_softwares.name, glpi_softwares.comment, IFNULL(glpi_entities.completename,'') AS ename, glpi_softwares.is_recursive,
		IFNULL(glpi_locations.completename,'') AS locations, glpi_softwares.groups_id_tech, glpi_softwares.users_id_tech, IFNULL(glpi_manufacturers.name,'') AS manufacture,
		IFNULL(softadd.moredescriptionfield,'') AS description1,
		 IFNULL(softadd.servicemanualurlfieldtwo,'') murl, IFNULL(softadd.technicaldescriptionurlfield,'') AS durl, IFNULL(glpi_groups.name,'') as group_name
		from glpi_softwares
		INNER JOIN glpi_entities ON glpi_softwares.entities_id=glpi_entities.id
		LEFT JOIN glpi_locations ON glpi_softwares.locations_id=glpi_locations.id
		LEFT JOIN glpi_manufacturers ON glpi_softwares.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON glpi_softwares.id=softadd.items_id
		LEFT JOIN glpi_groups ON glpi_softwares.groups_id_tech=glpi_groups.id
		WHERE glpi_softwares.id=%d`, id)
	err = r.db.SelectOne(&software, sql)
	return software, err
}
