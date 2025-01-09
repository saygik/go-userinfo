package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Integration struct {
	URL     string `json:"url"`
	Context struct {
		Action string `json:"action"`
	} `json:"context"`
}

type Action struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Integration Integration `json:"integration"`
}

type Attachment struct {
	Text    string   `json:"text"`
	Actions []Action `json:"actions"`
}

type MattermostCommand struct {
	UserName    string `json:"user_name" binding:"required"`
	ChannelId   string `json:"channel_id" binding:"required"`
	ChannelName string `json:"channel_name"`
	Command     string `json:"command"`
	TeamDomain  string `json:"team_domain"`
	TeamId      string `json:"team_id"`
	Text        string `json:"text"`
	Token       string `json:"token"`
	TriggerId   string `json:"trigger_id"`
	UserId      string `json:"user_id"`
}

func (h *Handler) MattermostGLPICommand(c *gin.Context) {
	// Read body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// Parse URL-encoded form data
	values, err := url.ParseQuery(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	cmd := MattermostCommand{
		//		UserName:    values.Get("user_name"),
		Command: values.Get("command"),
		Text:    values.Get("text"),
		Token:   values.Get("token"),
		// TriggerId:   values.Get("trigger_id"),
		// UserId:      values.Get("user_id"),
		// ChannelId:   values.Get("channel_id"),
		// ChannelName: values.Get("channel_name"),
		// TeamDomain:  values.Get("team_domain"),
		// TeamId:      values.Get("team_id"),
	}
	jsonData := `[
        {
            "pretext": "This is the attachment pretext.",
            "text": "This is the attachment text.",
            "actions": [
                {
                    "id": "message",
                    "name": "Ephemeral Message",
                    "integration": {
                        "url": "http://127.0.0.1:7357",
                        "context": {
                            "action": "do_something_ephemeral"
                        }
                    }
                },
                {
                    "id": "update",
                    "name": "Update",
                    "integration": {
                        "url": "http://127.0.0.1:7357",
                        "context": {
                            "action": "do_something_update"
                        }
                    }
                }
            ]
        }
    ]`

	var attachments []Attachment
	err = json.Unmarshal([]byte(jsonData), &attachments)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	if len(cmd.Command) > 2 {
		switch cmd.Command {
		case "/glpi":
			if cmd.Token != "tdjao9bz3pgjueowf9h7oa3mur" {
				c.JSON(http.StatusOK, gin.H{
					"response_type": "in_channel",
					"text":          fmt.Sprintf("Неверный токен"),
				})

				return
			}
			if len(cmd.Text) > 0 {

			}
			c.JSON(http.StatusOK, gin.H{
				"response_type": "in_channel",
				"text":          fmt.Sprintf("GLPI"),
				"username":      "bot-notificator",

				"attachments": attachments,
			})
			return
		default:
			c.JSON(http.StatusOK, gin.H{
				"response_type": "in_channel",
				"text":          fmt.Sprintf("Команда не распознана"),
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"response_type": "in_channel",
			"text":          fmt.Sprintf("Команда не распознана"),
		})
		return
	}

	// parts := strings.Fields(cmd.Text)
	// if len(parts) > 0 {
	// 	subcommand := parts[0]
	// 	switch subcommand {
	// 	case "subcommand1":
	// 		// Handle subcommand1
	// 		fmt.Sprintf("You triggered subcommand1")
	// 	case "subcommand2":
	// 		// Handle subcommand2
	// 		fmt.Sprintf("You triggered subcommand2")
	// 	default:
	// 		fmt.Sprintf("Unknown subcommand")
	// 	}
	// } else {
	// 	fmt.Sprintf("No subcommand provided")
	// }

}
