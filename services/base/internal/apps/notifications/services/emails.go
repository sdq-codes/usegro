package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/base/repositories"
	notificationModels "github.com/sdq-codes/usegro-api/internal/apps/notifications/models"
	"github.com/sdq-codes/usegro-api/internal/helper/email"
	"github.com/sdq-codes/usegro-api/internal/helper/queue"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"gorm.io/gorm"
)

func QueueEmail(ctx context.Context, tx *gorm.DB, rdb *redis.Client, emailNotification notificationModels.EmailNotification, jobRepository repositories.JobRepositoryInterface) error {
	emailNotificationBytes, err := json.Marshal(emailNotification)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("email job queue failed: %v", err))
		return err
	}
	job := &models.Job{
		ID:          uuid.New(),
		Queue:       "crm_queue_email",
		HandlerName: "SendEmail",
		Payload:     emailNotificationBytes,
	}

	newQueue := queue.NewQueue(job.Queue, jobRepository)
	if err = newQueue.Enqueue(ctx, tx, job); err != nil {
		return err
	}
	logger.Log.Info("email job queued")
	return nil
}

// SendEmail sends an email directly from the raw job payload bytes.
// The payload is already in the Redis queue message — no DB fetch needed.
func SendEmail(ctx context.Context, rawPayload json.RawMessage) error {
	var notification notificationModels.EmailNotification
	if err := json.Unmarshal(rawPayload, &notification); err != nil {
		logger.Log.Error(fmt.Sprintf("SendEmail: failed to unmarshal payload: %v", err))
		return err
	}

	mailRequest := email.MailRequest{
		From:    notification.FromEmail,
		To:      notification.ToEmails,
		Subject: notification.Subject,
	}
	mailRequest.Send(notification.Template, notification.Data)
	return nil
}
