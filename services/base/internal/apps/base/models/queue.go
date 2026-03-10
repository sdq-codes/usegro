package models

import (
	"github.com/google/uuid"
)

type QueuePayload struct {
	JobID uuid.UUID `json:"job_id"`
}
