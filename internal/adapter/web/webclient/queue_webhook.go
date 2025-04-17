package webclient

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) AddWebhook(data entity.WebhookPayload) error {
	r.WebhookQueue <- data
	return nil
}

func (r *Repository) Log() error {
	r.log.Info("Webhook added to queue")
	return nil
}
