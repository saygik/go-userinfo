package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatRegions(startdate string, enddate string) (tickets []entity.GLPIRegionsStats, err error) {
	sql := fmt.Sprintf(
		`SELECT count, org  FROM (
			SELECT count(glpi_tickets.id) AS COUNT, 'ИРЦ Минск' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИРЦ%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ2' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ2%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ3' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ3%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ4' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ4%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ5' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ5%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ6' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ6%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') a1
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
