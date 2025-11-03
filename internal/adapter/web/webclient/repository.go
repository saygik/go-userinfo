package webclient

import (
	"context"
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

func New(ctx context.Context, url string, log *logrus.Logger) *Repository {
	// Create a channel to act as the queue
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
	for {
		select {
		case <-r.ctx.Done():
			return
		case p := <-r.WebhookQueue:
			r.processOne(p) // без вложенной горутины
		}
	}
}
func (r *Repository) processOne(p entity.WebhookPayload) {
	backoff := time.Second
	const maxBackoff = 10 * time.Minute
	for retries := 0; retries < 5; retries++ {
		r.log.Infof("Sending webhook by HRP user %s in GLPI ticket %d", p.Data.FIO, p.Data.Id)
		if err := r.sendWebhook(p.Data, r.url); err == nil {
			r.log.Info("sended")
			return
		}
		// ожидание с возможностью отмены
		t := time.NewTimer(backoff)
		select {
		case <-r.ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
		if backoff *= 2; backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
	r.log.Infof("Max retries reached. Giving up on webhook: %s", p.WebhookId)
}

// func (r *Repository) ProcessWebhooks1() {
// 	for payload := range r.WebhookQueue {
// 		go func(p entity.WebhookPayload) {
// 			backoffTime := time.Second  // starting backoff time
// 			maxBackoffTime := time.Hour // maximum backoff time
// 			retries := 0
// 			maxRetries := 5

// 			for {
// 				r.log.Info("Sending webhook by HRP user " + p.Data.FIO + " in GLPI ticket " + strconv.Itoa(payload.Data.Id))
// 				err := r.sendWebhook(p.Data, r.url)
// 				if err == nil {
// 					r.log.Info("sended")
// 					break
// 				}
// 				r.log.Error("Error sending webhook:", err)
// 				//				r.log.Println("Error sending webhook:", err)

// 				retries++
// 				if retries >= maxRetries {
// 					r.log.Println("Max retries reached. Giving up on webhook:", p.WebhookId)
// 					break
// 				}

// 				time.Sleep(backoffTime)

// 				// Double the backoff time for the next iteration, capped at the max
// 				backoffTime *= 2
// 				r.log.Println(backoffTime)
// 				if backoffTime > maxBackoffTime {
// 					backoffTime = maxBackoffTime
// 				}
// 			}
// 		}(payload)
// 	}
// }
