package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatTickets() ([]entity.GLPITicketsStats, error) {
	tickets := []entity.GLPITicketsStats{}
	sql := `SELECT COUNT(id) AS count, TYPE AS type, YEAR (DATE) AS year, MONTH (DATE)  AS month FROM (SELECT * from glpi_tickets WHERE glpi_tickets.is_deleted=0) gt
		WHERE YEAR (DATE)>2020 GROUP BY MONTH (date) , YEAR (DATE), TYPE ORDER BY DATE`
	_, err := r.db.Select(&tickets, sql)
	return tickets, err
}
