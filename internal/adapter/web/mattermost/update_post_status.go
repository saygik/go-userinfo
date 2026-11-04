package mattermost

import (
	"context"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

// UpdatePostStatus добавляет статус к исходному сообщению Mattermost и, если
// задан цвет, меняет цвет полосы вложения. Повторное добавление того же
// статуса и цвета не выполняется.
func (r *Repository) UpdatePostStatus(postID, status, color string) error {
	if postID == "" || (status == "" && color == "") {
		return nil
	}

	post, _, err := r.client.GetPost(context.Background(), postID, "")
	if err != nil {
		return err
	}

	changed := false
	if status != "" && !strings.Contains(post.Message, status) {
		message := strings.TrimSpace(post.Message)
		if message != "" {
			message += "\n\n"
		}
		post.Message = message + "**" + status + "**"
		changed = true
	}
	if color != "" && setPostAttachmentColor(post, color) {
		changed = true
	}
	if !changed {
		return nil
	}

	_, _, err = r.client.UpdatePost(context.Background(), postID, post)
	return err
}

// setPostAttachmentColor меняет цвет полосы у вложений поста.
// Возвращает true, если хотя бы одно вложение изменилось.
func setPostAttachmentColor(post *model.Post, color string) bool {
	props := post.GetProps()
	if props == nil {
		return false
	}

	changed := false
	switch attachments := props["attachments"].(type) {
	case []*model.SlackAttachment:
		for _, attachment := range attachments {
			if attachment != nil && attachment.Color != color {
				attachment.Color = color
				changed = true
			}
		}
	case []any:
		for _, raw := range attachments {
			attachment, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if attachment["color"] != color {
				attachment["color"] = color
				changed = true
			}
		}
	case []map[string]any:
		for _, attachment := range attachments {
			if attachment["color"] != color {
				attachment["color"] = color
				changed = true
			}
		}
	}

	if changed {
		post.SetProps(props)
	}
	return changed
}
