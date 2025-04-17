package webclient

import (
	"context"
	"strconv"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	WebhookQueue chan entity.WebhookPayload
	url          string
	log          *logrus.Logger
	ctx          context.Context
}

func New(url string, log *logrus.Logger) *Repository {
	// Create a channel to act as the queue
	ctx := context.Background()
	//	log.Info("Starting webhook queue")
	webhookQueue := make(chan entity.WebhookPayload, 100) // Buffer size 100

	r := &Repository{
		WebhookQueue: webhookQueue,
		log:          log,
		url:          url,
		ctx:          ctx,
	}

	go r.ProcessWebhooks()
	return r
}

func (r *Repository) ProcessWebhooks() {
	for payload := range r.WebhookQueue {
		go func(p entity.WebhookPayload) {
			backoffTime := time.Second  // starting backoff time
			maxBackoffTime := time.Hour // maximum backoff time
			retries := 0
			maxRetries := 5

			for {
				r.log.Info("Sending webhook by HRP user " + p.Data.FIO + " in GLPI ticket " + strconv.Itoa(payload.Data.Id))
				err := r.sendWebhook(p.Data, r.url)
				if err == nil {
					r.log.Info("sended")
					break
				}
				r.log.Error("Error sending webhook:", err)
				//				r.log.Println("Error sending webhook:", err)

				retries++
				if retries >= maxRetries {
					r.log.Println("Max retries reached. Giving up on webhook:", p.WebhookId)
					break
				}

				time.Sleep(backoffTime)

				// Double the backoff time for the next iteration, capped at the max
				backoffTime *= 2
				r.log.Println(backoffTime)
				if backoffTime > maxBackoffTime {
					backoffTime = maxBackoffTime
				}
			}
		}(payload)
	}
}
