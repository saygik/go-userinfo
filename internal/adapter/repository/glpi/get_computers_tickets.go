package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

// func (r *Repository) GetComputersTickets(id int) (software map[string][]entity.IdName, err error) {
// 	// Первый запрос: базовая информация о ПО
// 	sql1 := fmt.Sprintf(`
//         SELECT
//             s.id,
//             s.name,
//             COALESCE(e.completename, '') AS ename,
//             '' AS login,
//             s.comment,
//             IFNULL(l.completename, '') AS locations,
//             IFNULL(m.name, '') AS manufacture,
//             IFNULL(softadd.moredescriptionfield, '') AS description1,
//             '' AS description2,
//             IFNULL(softadd.servicemanualurlfield, '') AS murl,
//             IFNULL(softadd.iconurlfield, '') AS durl,
//             s.is_recursive
//         FROM glpi_softwares s
//         INNER JOIN glpi_entities e ON s.entities_id = e.id
//         LEFT JOIN glpi_locations l ON s.locations_id = l.id
//         LEFT JOIN glpi_manufacturers m ON s.manufacturers_id = m.id
//         LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON s.id = softadd.items_id
//         WHERE s.id = %d AND s.is_deleted = 0`, id)

// 	err = r.db.SelectOne(&software, sql1)
// 	if err != nil {
// 		return software, err
// 	}

// 	// Второй запрос: группы техподдержки
// 	sql2 := fmt.Sprintf(`
//         SELECT
//             g.id,
//             g.name,
//             COALESCE(gmm.idmattermostfield, '') AS group_matt_channel,
//             CASE
//                 WHEN gmm.iduserinfofield REGEXP '^[0-9]+$'
//                 THEN CAST(gmm.iduserinfofield AS UNSIGNED)
//                 ELSE 0
//             END AS group_calendar
//         FROM glpi_groups_items gi
//         LEFT JOIN glpi_groups g ON gi.groups_id = g.id
//         LEFT JOIN glpi_plugin_fields_groupidmattermosts gmm ON g.id = gmm.items_id
//         WHERE gi.items_id = %d AND gi.itemtype = 'Software' AND gi.type = 2`, id)

// 	var groups []entity.SoftwareGroup
// 	_, err = r.db.Select(&groups, sql2)
// 	if err != nil {
// 		return software, err
// 	}

// 	software.Groups = groups

// 	return software, nil
// }

// 5.5 📦 Загружаем тикеты по компьютерам
func (r *Repository) GetComputersTickets() (ticketsByComp map[string][]entity.IdName, err error) {
	var tickets []struct {
		Computer string `db:"computer"`
		Id       int    `db:"id"`
		Name     string `db:"name"`
	}

	_, err = r.db.Select(&tickets, `
        SELECT c.name AS computer, ti.id, ti.name
        FROM glpi_items_tickets it
        LEFT JOIN glpi_tickets ti ON ti.id = it.tickets_id
        LEFT JOIN glpi_computers c ON c.id = it.items_id
        WHERE it.itemtype = 'Computer' AND ti.requesttypes_id != 13
    `)
	if err != nil {
		// логгировать, но не возвращать из функции ошибку, если тикеты не критичны
		return ticketsByComp, err
	} else {
		ticketsByComp = make(map[string][]entity.IdName, len(tickets))
		for _, t := range tickets {
			ticketsByComp[t.Computer] = append(ticketsByComp[t.Computer], entity.IdName{
				Id:   t.Id,
				Name: t.Name,
			})
		}
	}
	return ticketsByComp, nil
}
