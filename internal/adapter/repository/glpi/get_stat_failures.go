package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatFailures() ([]entity.GLPITicketsStats, error) {
	tickets := []entity.GLPITicketsStats{}
	sql := `SELECT COUNT(id) AS count, YEAR (DATE) AS year, MONTH (DATE)  AS month  FROM (SELECT glpi_tickets.id, glpi_tickets.date from glpi_tickets
			LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			INNER JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_tickets.id
			INNER JOIN ( select id from glpi_plugin_fields_failcategoryfielddropdowns WHERE id>4) gpf ON gpf.id=glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id
			WHERE  glpi_tickets.is_deleted<>TRUE ) d1
			WHERE YEAR (DATE)>2020
			GROUP BY MONTH (date) , YEAR (date)
			ORDER BY date`
	_, err := r.db.Select(&tickets, sql)
	return tickets, err
}
