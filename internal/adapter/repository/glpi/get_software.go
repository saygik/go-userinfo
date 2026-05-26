package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetSoftware(id int) (software entity.Software, err error) {
	// Первый запрос: базовая информация о ПО
	sql1 := fmt.Sprintf(`
        SELECT
            s.id,
            s.name,
            COALESCE(e.completename, '') AS ename,
            '' AS login,
            COALESCE(s.comment, '') AS comment,
            IFNULL(l.completename, '') AS locations,
            IFNULL(m.name, '') AS manufacture,
            IFNULL(softadd.moredescriptionfield, '') AS description1,
            '' AS description2,
            IFNULL(softadd.servicemanualurlfield, '') AS murl,
            IFNULL(softadd.iconurlfield, '') AS durl,
            s.is_recursive
        FROM glpi_softwares s
        INNER JOIN glpi_entities e ON s.entities_id = e.id
        LEFT JOIN glpi_locations l ON s.locations_id = l.id
        LEFT JOIN glpi_manufacturers m ON s.manufacturers_id = m.id
        LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON s.id = softadd.items_id
        WHERE s.id = %d AND s.is_deleted = 0`, id)

	err = r.db.SelectOne(&software, sql1)
	if err != nil {
		return software, err
	}

	// Второй запрос: группы техподдержки
	sql2 := fmt.Sprintf(`
        SELECT
            g.id,
            g.name,
            COALESCE(gmm.idmattermostfield, '') AS group_matt_channel,
            CASE
                WHEN gmm.iduserinfofield REGEXP '^[0-9]+$'
                THEN CAST(gmm.iduserinfofield AS UNSIGNED)
                ELSE 0
            END AS group_calendar
        FROM glpi_groups_items gi
        LEFT JOIN glpi_groups g ON gi.groups_id = g.id
        LEFT JOIN glpi_plugin_fields_groupidmattermosts gmm ON g.id = gmm.items_id
        WHERE gi.items_id = %d AND gi.itemtype = 'Software' AND gi.type = 2`, id)

	var groups []entity.SoftwareGroup
	_, err = r.db.Select(&groups, sql2)
	if err != nil {
		return software, err
	}

	software.Groups = groups

	return software, nil
}

func (r *Repository) GetSoftwareJournal(id int) (items []entity.SoftwareJournal, err error) {
	// Первый запрос: базовая информация о ПО
	sql1 := fmt.Sprintf(`SELECT
                CONCAT('C-', c.id) AS id,
                CONCAT('https://support.rw/front/change.form.php?id=', c.id) AS url,
                'изменение' AS doc_type,
                c.id AS item_id,
                c.name,
                c.content,
                c.date AS date_creation,
                'работа' AS request_type,
                '-' AS fail_category,
                0 AS fail_category_id                
                FROM glpi_changes c
                INNER JOIN glpi_changes_items ci
                ON ci.changes_id = c.id
                WHERE ci.itemtype = 'Software'
                AND c.is_deleted = 0
                AND ci.items_id = %d

                UNION ALL

                SELECT
                CONCAT('T-', t.id) AS id,
                CONCAT('https://support.rw/front/ticket.form.php?id=', t.id) AS url,
                'заявка' AS doc_type,
                t.id AS item_id,
                t.name,
                t.content,
                t.date AS date_creation,
                CASE
                    WHEN t.type = 1 THEN 'инцидент'
                    WHEN t.type = 2 THEN 'запрос'
                    ELSE 'Неизвестно'
                END AS request_type,
                IFNULL(fc.name, '-') AS fail_category,
                IFNULL(fc.id,0) AS fail_category_id
                FROM glpi_tickets t
                INNER JOIN glpi_items_tickets it
                ON it.tickets_id = t.id
                LEFT JOIN glpi_plugin_fields_ticketfailures tf
                ON tf.items_id = t.id
                LEFT JOIN glpi_plugin_fields_failcategoryfielddropdowns fc
                ON fc.id = tf.plugin_fields_failcategoryfielddropdowns_id
                WHERE it.itemtype = 'Software'
                AND t.is_deleted = 0
                AND it.items_id = %d
                ORDER BY date_creation DESC`, id, id)

	_, err = r.db.Select(&items, sql1)
	if err != nil {
		return nil, err
	}
	return items, nil
}
