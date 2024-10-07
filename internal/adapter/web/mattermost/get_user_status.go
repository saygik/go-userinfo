package mattermost

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

func (r *Repository) GetUserStatus(id string) (status *model.Status, err error) {
	status, _, err = r.client.GetUserStatus(id, "")
	return status, err
}
