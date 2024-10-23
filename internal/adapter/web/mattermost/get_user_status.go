package mattermost

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
)

func (r *Repository) GetUserStatus(id string) (status *model.Status, err error) {
	status, _, err = r.client.GetUserStatus(context.Background(), id, "")
	return status, err
}
