package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetAllSoftwares() (softwares []entity.Software, err error) {
	// Первый запрос: базовая информация о ПО (без групп)
	sql1 := `
        SELECT
            s.id, s.name, COALESCE(s.comment,'') AS comment , COALESCE(e.completename,'') AS ename,
            s.is_recursive, IFNULL(l.completename,'') AS locations,
            IFNULL(m.name,'') AS manufacture,
            IFNULL(softadd.moredescriptionfield,'') AS description1,
            IFNULL(softadd.servicemanualurlfield,'') AS murl,
            IFNULL(softadd.iconurlfield,'') AS durl
        FROM glpi_softwares s
        INNER JOIN glpi_entities e ON s.entities_id = e.id
        LEFT JOIN glpi_locations l ON s.locations_id = l.id
        LEFT JOIN glpi_manufacturers m ON s.manufacturers_id = m.id
        LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON s.id = softadd.items_id
         INNER JOIN glpi_groups_items gi ON s.id = gi.items_id AND gi.itemtype = 'Software' AND gi.type = 2
        WHERE s.is_deleted = 0
        ORDER BY s.name`

	_, err = r.db.Select(&softwares, sql1)
	if err != nil {
		return nil, err
	}

	// Второй запрос: все группы для всех ПО (с подзапросом)
	sql2 := `
        SELECT
            s.id AS software_id,
            g.id,
            g.name,
            COALESCE(gmm.idmattermostfield, '') AS group_matt_channel,
            CASE
                WHEN gmm.iduserinfofield REGEXP '^[0-9]+$'
                THEN CAST(gmm.iduserinfofield AS UNSIGNED)
                ELSE 0
            END AS group_calendar
        FROM glpi_softwares s
        INNER JOIN glpi_groups_items gi ON s.id = gi.items_id
            AND gi.itemtype = 'Software' AND gi.type = 2
        LEFT JOIN glpi_groups g ON gi.groups_id = g.id
        LEFT JOIN glpi_plugin_fields_groupidmattermosts gmm ON g.id = gmm.items_id
        WHERE s.is_deleted = 0
        ORDER BY s.id, g.name`

	type SoftwareGroupWithID struct {
		SoftwareID int64 `db:"software_id"`
		entity.SoftwareGroup
	}

	var allGroups []SoftwareGroupWithID
	_, err = r.db.Select(&allGroups, sql2)
	if err != nil {
		return nil, err
	}

	// Группируем группы по software_id
	groupMap := make(map[int64][]entity.SoftwareGroup)
	for _, sg := range allGroups {
		groupMap[sg.SoftwareID] = append(groupMap[sg.SoftwareID], sg.SoftwareGroup)
	}

	// Присваиваем группы к ПО
	for i := range softwares {
		softwares[i].Groups = groupMap[softwares[i].Id]
	}

	return softwares, nil
}
