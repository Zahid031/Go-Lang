package models

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID        int64           `json:"id"`
	JobType   string          `json:"job_type"`
	Payload   json.RawMessage `json:"payload"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}