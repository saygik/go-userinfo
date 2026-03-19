package glpi

import (
	"strconv"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetAllSoftwares() (softwares []entity.Software, err error) {
	sql := `
       SELECT
            s.id, s.name, s.comment, IFNULL(e.completename,'') AS ename,
            s.is_recursive, IFNULL(l.completename,'') AS locations,
            s.users_id_tech, IFNULL(m.name,'') AS manufacture,
            IFNULL(softadd.moredescriptionfield,'') AS description1,
            IFNULL(softadd.servicemanualurlfield,'') AS murl,
            IFNULL(softadd.iconurlfield,'') AS durl,
            COALESCE(GROUP_CONCAT(DISTINCT gi.groups_id), '') AS groups_id_tech_s,
            COALESCE(GROUP_CONCAT(DISTINCT g.name SEPARATOR ', '), '') AS group_names_s
        FROM glpi_softwares s
        INNER JOIN glpi_entities e ON s.entities_id = e.id
        LEFT JOIN glpi_locations l ON s.locations_id = l.id
        LEFT JOIN glpi_manufacturers m ON s.manufacturers_id = m.id
        LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON s.id = softadd.items_id
        LEFT JOIN glpi_groups_items gi ON s.id = gi.items_id
            AND gi.itemtype = 'Software' AND gi.type = 2
        LEFT JOIN glpi_groups g ON gi.groups_id = g.id
        WHERE s.is_deleted = 0
        GROUP BY s.id, s.name, s.comment, e.completename, s.is_recursive,
                 l.completename, s.users_id_tech, m.name, softadd.moredescriptionfield,
                 softadd.servicemanualurlfield, softadd.iconurlfield
        ORDER BY s.name
    `

	_, err = r.db.Select(&softwares, sql)

	// 🔥 Парсим строковые поля в массивы
	for i := range softwares {
		softwares[i].Groups_id_tech = parseGroupIDs(softwares[i].Groups_id_tech_s)
		softwares[i].GroupNames = parseGroupNames(softwares[i].GroupNames_s)
	}

	return softwares, err
}

// Парсит "1,2,5" → []int64{1,2,5}
func parseGroupIDs(idsStr string) []int64 {
	if idsStr == "" || idsStr == "NULL" {
		return []int64{}
	}

	idStrings := strings.Split(idsStr, ",")
	groupIDs := make([]int64, 0, len(idStrings))

	for _, idStr := range idStrings {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			groupIDs = append(groupIDs, id)
		}
	}
	return groupIDs
}

// Парсит "Group1, Group2" → []string{"Group1", "Group2"}
func parseGroupNames(namesStr string) []string {
	if namesStr == "" || namesStr == "NULL" {
		return []string{}
	}

	return strings.FieldsFunc(namesStr, func(r rune) bool {
		return r == ','
	})
}
