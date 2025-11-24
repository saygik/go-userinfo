package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicketNetworkEquipment(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
		SELECT  glpi_networkequipments.name AS 'name',
		CONCAT (IFNULL(glpi_networkequipmenttypes.name,""), " ",IFNULL(glpi_manufacturers.name,""), " ",IFNULL(glpi_networkequipmentmodels.name,"") )  AS fullname,
		IFNULL(glpi_groups.name, "") AS 'group',
		IFNULL(glpi_locations.name, "") AS place
		FROM (SELECT * FROM glpi_items_tickets WHERE itemtype='NetworkEquipment' and tickets_id=%s) it
		LEFT JOIN glpi_networkequipments    ON it.items_id=glpi_networkequipments.id
		LEFT JOIN glpi_groups ON glpi_networkequipments.groups_id_tech=glpi_groups.id
		LEFT JOIN glpi_locations ON glpi_networkequipments.locations_id=glpi_locations.id
		LEFT JOIN glpi_networkequipmenttypes ON glpi_networkequipments.networkequipmenttypes_id=glpi_networkequipmenttypes.id
		LEFT JOIN glpi_manufacturers ON glpi_networkequipments.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_networkequipmentmodels ON glpi_networkequipments.networkequipmentmodels_id=glpi_networkequipmentmodels.id
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (r *Repository) GetTicketSoft(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
		SELECT  glpi_softwares.name AS 'name',
		IFNULL(glpi_groups.name, "") AS 'group',
		'' AS fullname, IFNULL(glpi_locations.name, "") as place
		FROM (SELECT * FROM glpi_items_tickets WHERE itemtype='Software' and tickets_id=%s) it
		LEFT JOIN glpi_softwares   ON it.items_id=glpi_softwares.id
		LEFT JOIN glpi_groups ON glpi_softwares.groups_id_tech=glpi_groups.id
		LEFT JOIN glpi_locations ON glpi_softwares.locations_id=glpi_locations.id
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (r *Repository) GetTicketServers(ticketID string) (obj []entity.GLPI_Object, err error) {
	var proc = fmt.Sprintf(`
		SELECT  glpi_computers.name AS 'name',
		IFNULL(glpi_groups.name, "") AS 'group',
		CONCAT (IFNULL(glpi_computertypes.name,""), " ",IFNULL(glpi_manufacturers.name,""), " ",IFNULL(glpi_computermodels.name,""))  AS fullname,
		IFNULL(glpi_locations.name, "") as place
		FROM (SELECT * FROM glpi_items_tickets WHERE itemtype='Computer' and tickets_id=%s) it
		LEFT JOIN glpi_computers   ON it.items_id=glpi_computers.id
		LEFT JOIN glpi_groups ON glpi_computers.groups_id_tech=glpi_groups.id
		LEFT JOIN glpi_locations ON glpi_computers.locations_id=glpi_locations.id
		LEFT JOIN glpi_computertypes ON glpi_computers.computertypes_id=glpi_computertypes.id
		LEFT JOIN glpi_manufacturers ON glpi_computers.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_computermodels ON glpi_computers.computermodels_id=glpi_computermodels.id
	`, ticketID)
	_, err = r.db.Select(&obj, proc)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
