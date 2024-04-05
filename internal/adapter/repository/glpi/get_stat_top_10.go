package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatTop10Performers(startdate string, enddate string) (tickets []entity.GLPIStatsTop10, err error) {
	sql := fmt.Sprintf(
		`SELECT * FROM ( SELECT  gg.name, TRIM(CONCAT(ifnull(NULLIF(gg.realname, ''), ''),' ', ifnull(NULLIF(gg.firstname, ''),''))) AS completename,ge.completename AS company,
		 count(gt.id) AS count FROM (SELECT * from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND ((glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') OR (glpi_tickets.solvedate>= '%[1]s' AND glpi_tickets.solvedate<= '%[2]s'))) gt
		 inner JOIN (SELECT id, tickets_id, users_id FROM glpi_tickets_users WHERE TYPE=2) ggt ON ggt.tickets_id= gt.id
		 inner JOIN glpi_users gg ON gg.id = ggt.users_id
		 inner JOIN glpi_entities ge ON ge.id = gg.entities_id
		 GROUP BY gg.name, ge.completename ) tab
		 ORDER BY COUNT  DESC LIMIT 10
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}

func (r *Repository) GetStatTop10Iniciators(startdate string, enddate string) (tickets []entity.GLPIStatsTop10, err error) {
	sql := fmt.Sprintf(
		`SELECT * FROM ( SELECT  gg.name, TRIM(CONCAT(ifnull(NULLIF(gg.realname, ''), ''),' ', ifnull(NULLIF(gg.firstname, ''),''))) AS completename,ge.completename AS company,
		 count(gt.id) AS count FROM (SELECT * from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') gt
		 inner JOIN (SELECT id, tickets_id, users_id FROM glpi_tickets_users WHERE TYPE=1) ggt ON ggt.tickets_id= gt.id
		 inner JOIN glpi_users gg ON gg.id = ggt.users_id
		 inner JOIN glpi_entities ge ON ge.id = gg.entities_id
		 GROUP BY gg.name, ge.completename ) tab
		 ORDER BY COUNT  DESC LIMIT 10
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}

func (r *Repository) GetStatTop10Groups(startdate string, enddate string) (tickets []entity.GLPIStatsTop10, err error) {
	sql := fmt.Sprintf(
		`SELECT * FROM ( SELECT gg.name, ifnull(gg.comment,'') AS completename, ge.completename AS company,  count(gt.id) AS count
		FROM (SELECT * from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND ((glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s') OR (glpi_tickets.solvedate>= '%[1]s' AND glpi_tickets.solvedate<= '%[2]s'))) gt
		inner JOIN (SELECT id, tickets_id, groups_id FROM glpi_groups_tickets WHERE TYPE=2) ggt ON ggt.tickets_id= gt.id
		 inner JOIN glpi_groups gg ON gg.id = ggt.groups_id
		 inner JOIN glpi_entities ge ON ge.id = gg.entities_id
		 GROUP BY gg.name, ge.completename ) tab
		 ORDER BY count DESC LIMIT 10
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
