package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatPeriodRequestTypes(startdate string, enddate string) (tickets []entity.GLPIStatsTop10, err error) {
	sql := fmt.Sprintf(
		`SELECT (case  requesttypes_id
     		when 1 then 'Техподдержка'
     		when 2 then 'E-Mail'
     		when 3 then 'Телефон'
     		when 4 then 'Уснто'
     		when 5 then 'Письмо'
     		when 6 then 'Другой'
     		when 7 then 'WEB форма'
     		when 8 then 'SAP HRP'
     		when 9 then 'Mattermost'
     		ELSE 'Техподдержка'
			END)  AS name,
			COUNT(id) AS count from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND glpi_tickets.date>= '%[1]s' AND glpi_tickets.date<= '%[2]s'
			GROUP BY requesttypes_id
		 `, startdate, enddate)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
