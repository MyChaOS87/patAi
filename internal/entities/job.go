package entities

import (
	"github.com/google/uuid"
)

type EvaluationJobStatus int

const (
	EvaluationJobStatusPending EvaluationJobStatus = iota
	EvaluationJobStatusFinished
	EvaluationJobStatusFailed
)

type EvaluationJob struct {
	ID                  uuid.UUID
	EvaluationJobStatus EvaluationJobStatus
	PatentContent       string
	Value               int
	OwnerID             string
}
