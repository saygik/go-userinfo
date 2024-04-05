package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatTicketsDays(startdate string, enddate string) (tickets []entity.GLPITicketsStats, err error) {
	sql := fmt.Sprintf(
		`SELECT COUNT(id) AS count, TYPE AS type, YEAR (DATE) AS year, MONTH (DATE)  AS month, DAY (DATE) AS day
		 FROM (SELECT * from glpi_tickets WHERE glpi_tickets.is_deleted=0 AND date>='%[1]s' AND date <='%[2]s') gt
		 GROUP BY DAY (DATE), MONTH (date) , YEAR (DATE), TYPE ORDER BY DATE`, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
