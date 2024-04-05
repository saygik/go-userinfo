package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetStatPeriodRegionDayCounts(startdate string, enddate string, maxday int) (tickets []entity.RegionsDayStats, err error) {
	sql := fmt.Sprintf(
		`SELECT DATE AS 'Day', org , ifnull(COUNT, 0) AS count FROM (
    SELECT 1 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 1 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 1 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 1 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 1 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 1 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 2 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 2 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 2 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 2 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 2 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 2 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 3 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 3 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 3 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 3 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 3 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 3 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 4 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 4 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 4 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 4 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 4 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 4 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 5 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 5 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 5 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 5 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 5 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 5 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 6 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 6 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 6 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 6 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 6 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 6 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 7 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 7 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 7 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 7 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 7 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 7 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 8 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 8 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 8 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 8 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 8 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 8 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 9 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 9 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 9 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 9 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 9 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 9 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 10 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 10 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 10 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 10 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 10 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 10 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 11 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 11 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 11 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 11 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 11 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 11 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 12 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 12 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 12 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 12 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 12 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 12 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 13 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 13 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 13 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 13 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 13 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 13 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 14 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 14 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 14 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 14 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 14 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 14 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 15 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 15 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 15 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 15 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 15 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 15 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 16 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 16 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 16 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 16 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 16 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 16 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 17 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 17 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 17 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 17 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 17 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 17 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 18 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 18 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 18 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 18 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 18 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 18 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 19 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 19 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 19 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 19 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 19 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 19 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 20 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 20 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 20 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 20 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 20 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 20 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 21 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 21 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 21 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 21 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 21 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 21 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 22 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 22 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 22 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 22 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 22 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 22 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 23 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 23 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 23 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 23 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 23 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 23 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 24 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 24 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 24 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 24 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 24 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 24 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 25 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 25 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 25 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 25 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 25 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 25 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 26 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 26 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 26 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 26 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 26 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 26 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 27 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 27 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 27 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 27 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 27 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 27 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 28 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 28 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 28 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 28 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 28 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 28 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 29 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 29 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 29 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 29 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 29 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 29 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 30 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 30 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 30 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 30 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 30 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 30 AS DATE, 'ИВЦ6' AS org UNION ALL
    SELECT 31 AS DATE, 'ИРЦ Минск' AS org UNION ALL
    SELECT 31 AS DATE, 'ИВЦ2' AS org UNION ALL
    SELECT 31 AS DATE, 'ИВЦ3' AS org UNION ALL
    SELECT 31 AS DATE, 'ИВЦ4' AS org UNION ALL
    SELECT 31 AS DATE, 'ИВЦ5' AS org UNION ALL
    SELECT 31, 'ИВЦ6' AS org) AS MonthDate
LEFT JOIN (SELECT COUNT(id) AS count, YEAR (DATE) AS year, MONTH (DATE)  AS month, DAY (DATE) AS DAY, name
FROM ( SELECT glpi_tickets.id,DATE, TRIM(SUBSTRING_INDEX(SUBSTRING(completename, 7), '>',1)) AS name from glpi_tickets
LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
WHERE glpi_tickets.entities_id>0 AND glpi_tickets.entities_id!=519 AND glpi_tickets.is_deleted=0 AND date>='%[1]s' AND date <='%[2]s'
) gt GROUP BY DAY (DATE), MONTH (date) , YEAR (DATE), name ORDER BY DATE) T1 ON MonthDate.Date = T1.DAY  AND MonthDate.org = T1.name
WHERE MonthDate.DATE <= %[3]d
ORDER BY Date`, startdate, enddate, maxday)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
