package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/notifications/services"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"gorm.io/gorm"
)

type WorkerService struct {
	rdb       *redis.Client
	db        *gorm.DB
	queues    []string
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewWorkerService(rdb *redis.Client, tx *gorm.DB, queues []string) *WorkerService {
	ctx, cancel := context.WithCancel(context.Background())
	logger.Log.Info("Starting workers...")
	return &WorkerService{
		rdb:       rdb,
		db:        tx,
		queues:    queues,
		ctx:       ctx,
		cancelCtx: cancel,
	}
}

func (w *WorkerService) Start() {
	for _, queue := range w.queues {
		go w.listenQueue(queue)
	}
}

func (w *WorkerService) Stop() {
	w.cancelCtx()
}

func (w *WorkerService) listenQueue(queue string) {
	fmt.Printf("Listening to queue: %s\n", queue)

	for {
		select {
		case <-w.ctx.Done():
			fmt.Printf("Stopped listening to %s\n", queue)
			return
		default:
			result, err := w.rdb.BRPop(w.ctx, 5*time.Second, queue).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				fmt.Printf("Error reading from queue %s: %v\n", queue, err)
				continue
			}

			job := &models.Job{}
			if err = json.Unmarshal([]byte(result[1]), job); err != nil {
				logger.Log.Error(fmt.Sprintf("Error unmarshalling job: %v", err))
				continue
			}
			go w.processJob(queue, job)
		}
	}
}

func (w *WorkerService) processJob(queue string, payload *models.Job) {
	switch payload.HandlerName {
	case "SendEmail":
		w.handleEmail(payload)
	default:
		fmt.Printf("Unknown handler: %s\n", payload.HandlerName)
	}
}

func (w *WorkerService) handleEmail(payload *models.Job) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// payload.Payload already contains the full EmailNotification — no DB fetch needed.
	if err := services.SendEmail(ctx, json.RawMessage(payload.Payload)); err != nil {
		logger.Log.Error(fmt.Sprintf("handleEmail: failed to send email for job %s: %v", payload.ID, err))
	}
}
