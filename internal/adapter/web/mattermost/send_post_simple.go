package mattermost

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
)

func (r *Repository) SendPostSimple(
	channelId string,
	message string,
) (err error) {
	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
	}
	if _, _, err := r.client.CreatePost(context.Background(), post); err != nil {
		return err
	}
	return nil
}

func (r *Repository) ReplyPost(
	channelId string,
	postId string,
	message string,
) (err error) {
	post := &model.Post{
		ChannelId: channelId,
		RootId:    postId,
		Message:   message,
	}
	if _, _, err := r.client.CreatePost(context.Background(), post); err != nil {
		return err
	}
	return nil
}
