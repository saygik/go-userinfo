package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicketNetworkEquipment(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
			SELECT
				ne.name AS name,
				CONCAT(
					IFNULL(net.name, ''), ' ',
					IFNULL(m.name, ''), ' ',
					IFNULL(nem.name, '')
				) AS fullname,
				COALESCE(GROUP_CONCAT(DISTINCT g.name SEPARATOR ', '), '') AS 'group',
				IFNULL(l.name, '') AS place
			FROM (
				SELECT * FROM glpi_items_tickets
				WHERE itemtype = 'NetworkEquipment' AND tickets_id = %s
			) it
			LEFT JOIN glpi_networkequipments ne ON it.items_id = ne.id
			LEFT JOIN glpi_groups_items gi ON ne.id = gi.items_id
				AND gi.itemtype = 'NetworkEquipment'
				AND gi.type = 2
			LEFT JOIN glpi_groups g ON gi.groups_id = g.id
			LEFT JOIN glpi_locations l ON ne.locations_id = l.id
			LEFT JOIN glpi_networkequipmenttypes net ON ne.networkequipmenttypes_id = net.id
			LEFT JOIN glpi_manufacturers m ON ne.manufacturers_id = m.id
			LEFT JOIN glpi_networkequipmentmodels nem ON ne.networkequipmentmodels_id = nem.id
			WHERE ne.is_deleted = 0
			GROUP BY ne.id, ne.name, l.name, net.name, m.name, nem.name
			ORDER BY ne.name
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (r *Repository) GetTicketSoft(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
		SELECT
			s.name AS name,
			COALESCE(GROUP_CONCAT(DISTINCT g.name SEPARATOR ', '), '') AS 'group',
			'' AS fullname,
			IFNULL(l.name, '') AS place
		FROM (
			SELECT * FROM glpi_items_tickets
			WHERE itemtype = 'Software' AND tickets_id = %s
		) it
		LEFT JOIN glpi_softwares s ON it.items_id = s.id
		LEFT JOIN glpi_groups_items gi ON s.id = gi.items_id
			AND gi.itemtype = 'Software'
			AND gi.type = 2
		LEFT JOIN glpi_groups g ON gi.groups_id = g.id
		LEFT JOIN glpi_locations l ON s.locations_id = l.id
		WHERE s.is_deleted = 0
		GROUP BY s.id, s.name, l.name
		ORDER BY s.name
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (r *Repository) GetTicketServers(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
		SELECT
			c.name AS name,
			COALESCE(GROUP_CONCAT(DISTINCT g.name SEPARATOR ', '), '') AS 'group',
			CONCAT(
				IFNULL(ct.name, ''), ' ',
				IFNULL(m.name, ''), ' ',
				IFNULL(cm.name, '')
			) AS fullname,
			IFNULL(l.name, '') AS place
		FROM (
			SELECT * FROM glpi_items_tickets
			WHERE itemtype = 'Computer' AND tickets_id = %s
		) it
		LEFT JOIN glpi_computers c ON it.items_id = c.id
		LEFT JOIN glpi_groups_items gi ON c.id = gi.items_id
			AND gi.itemtype = 'Computer'
			AND gi.type = 2
		LEFT JOIN glpi_groups g ON gi.groups_id = g.id
		LEFT JOIN glpi_locations l ON c.locations_id = l.id
		LEFT JOIN glpi_computertypes ct ON c.computertypes_id = ct.id
		LEFT JOIN glpi_manufacturers m ON c.manufacturers_id = m.id
		LEFT JOIN glpi_computermodels cm ON c.computermodels_id = cm.id
		WHERE c.is_deleted = 0
		GROUP BY c.id, c.name, l.name, ct.name, m.name, cm.name
		ORDER BY c.name
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
