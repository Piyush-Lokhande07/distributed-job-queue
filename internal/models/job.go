package models

import "time"

const (

	StateQueued = "QUEUED"
	StateInProgress = "IN_PROGRESS"
	StateCompleted = "COMPLETED"
	StateFailed = "FAILED"
)

type Job struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Retries    int       `json:"retries"`
	MaxRetries int       `json:"max_retries"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
