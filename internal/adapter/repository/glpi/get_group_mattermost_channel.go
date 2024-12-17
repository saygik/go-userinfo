package glpi

func (r *Repository) GetGroupMattermostChannel(id int) (channelId string, calId int, err error) {

	//err = r.db.QueryRow(`SELECT glpi_plugin_fields_groupidmattermosts.idmattermostfield from glpi_plugin_fields_groupidmattermosts WHERE items_id=?`, id).Scan(&channelId)
	err = r.db.QueryRow(`SELECT glpi_plugin_fields_groupidmattermosts.idmattermostfield, glpi_plugin_fields_groupidmattermosts.iduserinfofield from glpi_plugin_fields_groupidmattermosts WHERE items_id=?`, id).Scan(&channelId, &calId)
	return channelId, calId, err

}
