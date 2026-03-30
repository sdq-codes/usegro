package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Queue        string    `json:"queue" gorm:"not null"`
	HandlerName  string    `json:"handler_name" gorm:"not null"`
	Payload      []byte    `json:"payload" gorm:"type:jsonb;not null"`
	MaxAttempts  int       `json:"max_attempts" gorm:"not null"`
	Delay        int       `json:"delay"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null"`
	FunctionName string    `json:"function_name" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	FailedJob []FailedJob `json:"failed_job" gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}

type FailedJob struct {
	ID       int             `json:"id" gorm:"primaryKey;autoIncrement"`
	JobID    uuid.UUID       `json:"job_id" gorm:"type:uuid;index;not null"`
	Queue    string          `json:"queue" gorm:"not null"`
	Payload  json.RawMessage `json:"payload" gorm:"type:jsonb"`
	Error    string          `json:"error" gorm:"type:text"`
	FailedAt time.Time       `json:"failed_at" gorm:"autoCreateTime"`

	Job Job `gorm:"foreignKey:JobID"`
}
