package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatPeriodOrgTreemap(startdate string, enddate string) (tickets []entity.TreemapData, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_entities.name AS x, count(glpi_tickets.id) AS y  FROM glpi_tickets
		INNER JOIN glpi_entities ON glpi_tickets.entities_id=glpi_entities.id
		WHERE glpi_tickets.is_deleted=0  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
		GROUP BY x ORDER BY y desc
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
