package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetAllRMSSoftwares(org string) (softwares []entity.ComputerRMS, err error) {
	// Первый запрос: базовая информация о ПО (без групп)
	sql1 := fmt.Sprintf(`
            SELECT  glpi_computers.name AS computer ,glpi_softwareversions.name AS version, glpi_items_softwareversions.date_install AS softdate 
            FROM glpi_softwareversions INNER JOIN glpi_items_softwareversions ON (glpi_items_softwareversions.softwareversions_id = glpi_softwareversions.id) 
            INNER JOIN glpi_computers ON (glpi_computers.id = glpi_items_softwareversions.items_id AND glpi_items_softwareversions.itemtype = 'Computer') 
            LEFT JOIN glpi_entities ON glpi_entities.id = glpi_computers.entities_id
            WHERE glpi_softwareversions.softwares_id = '3815' AND glpi_items_softwareversions.is_deleted = '0' AND (true) AND glpi_computers.is_deleted = '0' AND glpi_computers.is_template = '0'
            AND glpi_entities.completename LIKE '%s'`, org+"%")

	_, err = r.db.Select(&softwares, sql1)
	if err != nil {
		return nil, err
	}

	return softwares, nil
}
