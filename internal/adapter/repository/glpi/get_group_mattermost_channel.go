package glpi

func (r *Repository) GetGroupMattermostChannel(id int) (name string, channelId string, calId int, err error) {

	//err = r.db.QueryRow(`SELECT glpi_plugin_fields_groupidmattermosts.idmattermostfield from glpi_plugin_fields_groupidmattermosts WHERE items_id=?`, id).Scan(&channelId)
	err = r.db.QueryRow(`SELECT name, glpi_plugin_fields_groupidmattermosts.idmattermostfield, glpi_plugin_fields_groupidmattermosts.iduserinfofield
	from glpi_plugin_fields_groupidmattermosts INNER JOIN glpi_groups ON glpi_groups.id=glpi_plugin_fields_groupidmattermosts.items_id
	WHERE items_id=?`, id).Scan(&name, &channelId, &calId)
	return name, channelId, calId, err

}
