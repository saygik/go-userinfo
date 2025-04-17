package entity

type WebhookPayload struct {
	WebhookId string  `json:"webhookId"`
	Data      HRPUser `json:"data"`
}
