package glpi

func (r *Repository) GetGroupMattermostChannel(id int) (channelId string, err error) {

	err = r.db.QueryRow(`SELECT glpi_plugin_fields_groupidmattermosts.idmattermostfield from glpi_plugin_fields_groupidmattermosts WHERE items_id=?`, id).Scan(&channelId)

	return channelId, err

}
