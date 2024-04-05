package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatPeriodTicketsCounts(startdate string, enddate string) (tickets []entity.GLPIStatsCounts, err error) {
	sql := fmt.Sprintf(
		`SELECT (select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') AS t1,
		(select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND TYPE=1 AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') AS t1_1 ,
		(select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND TYPE=2 AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') AS t1_2,
		(select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND glpi_tickets.solvedate>= '%[1]s' AND glpi_tickets.solvedate<= '%[2]s') AS t2,
		(select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND TYPE=1 AND glpi_tickets.solvedate>= '%[1]s' AND glpi_tickets.solvedate<= '%[2]s') AS t2_1,
		(select COUNT(id) from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND TYPE=2 AND glpi_tickets.solvedate>= '%[1]s' AND glpi_tickets.solvedate<= '%[2]s') AS t2_2,
		(SELECT count(glpi_tickets.id) FROM glpi_tickets WHERE  glpi_tickets.is_deleted<>TRUE AND ( glpi_tickets.status = 4 OR glpi_tickets.status =3 OR glpi_tickets.status =2)) AS  t3,
		(SELECT count(glpi_tickets.id) FROM glpi_tickets WHERE  glpi_tickets.is_deleted<>TRUE AND ( glpi_tickets.status = 2 OR glpi_tickets.status =3)) AS  t3_1,
		(SELECT count(glpi_tickets.id) FROM glpi_tickets WHERE  glpi_tickets.is_deleted<>TRUE AND ( glpi_tickets.status = 3 )) AS  t3_2
		FROM glpi_tickets LIMIT 1
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
