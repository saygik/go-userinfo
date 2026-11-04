package mattermost

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestSetPostAttachmentColorFromJSON(t *testing.T) {
	post := &model.Post{}
	post.SetProps(model.StringInterface{
		"attachments": []any{
			map[string]any{"color": "#FF2200"},
		},
	})

	if !setPostAttachmentColor(post, "#808080") {
		t.Fatal("ожидалось изменение цвета вложения")
	}

	attachments := post.GetProps()["attachments"].([]any)
	color := attachments[0].(map[string]any)["color"]
	if color != "#808080" {
		t.Fatalf("цвет не изменился: %v", color)
	}
	if setPostAttachmentColor(post, "#808080") {
		t.Fatal("повторная установка того же цвета не должна менять пост")
	}
}

func TestSetPostAttachmentColorSlackAttachments(t *testing.T) {
	post := &model.Post{}
	post.SetProps(model.StringInterface{
		"attachments": []*model.SlackAttachment{{Color: "#FF2200"}},
	})

	if !setPostAttachmentColor(post, "#808080") {
		t.Fatal("ожидалось изменение цвета вложения")
	}

	attachments := post.GetProps()["attachments"].([]*model.SlackAttachment)
	if attachments[0].Color != "#808080" {
		t.Fatalf("цвет не изменился: %s", attachments[0].Color)
	}
}

func TestSetPostAttachmentColorWithoutAttachments(t *testing.T) {
	post := &model.Post{}
	if setPostAttachmentColor(post, "#808080") {
		t.Fatal("без вложений изменение цвета невозможно")
	}
}
