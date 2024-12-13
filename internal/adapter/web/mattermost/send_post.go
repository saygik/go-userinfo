package mattermost

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
)

func (r *Repository) SendPost(
	channelId string,
	name string,
	text string,
	title string,
	titleLink string,
	footer string,
	requestedAck bool,
) (err error) {
	post := &model.Post{
		ChannelId: channelId,
		Metadata: &model.PostMetadata{
			Priority: &model.PostPriority{
				Priority:     model.NewPointer("standard"), // Options: "standard", "important", "urgent"
				RequestedAck: model.NewPointer(requestedAck),
			},
		}}

	post.SetProps(map[string]interface{}{

		"attachments": []*model.SlackAttachment{
			{

				AuthorName: name,
				Text:       text,
				Color:      "#FFA500",
				Title:      title,
				TitleLink:  titleLink,
				Footer:     footer,
			},
		},
	})

	if _, _, err := r.client.CreatePost(context.Background(), post); err != nil {
		return err
	}
	return nil
}
