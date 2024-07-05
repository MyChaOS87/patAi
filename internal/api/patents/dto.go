package patents

import "github.com/MyChaOS87/patAi/internal/entities"

const (
	dtoStatusPending  = "pending"
	dtoStatusFinished = "finished"
	dtoStatusFailed   = "failed"
	dtoStatusUnknown  = "unknown"
)

type JobDTO struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Value  *int   `json:"value,omitempty"`
}

func JobToDTO(job entities.EvaluationJob) JobDTO {
	dto := JobDTO{
		ID:    job.ID.String(),
		Value: nil,
	}

	switch job.EvaluationJobStatus {
	case entities.EvaluationJobStatusPending:
		dto.Status = dtoStatusPending
	case entities.EvaluationJobStatusFinished:
		dto.Status = dtoStatusFinished
		dto.Value = &job.Value
	case entities.EvaluationJobStatusFailed:
		dto.Status = dtoStatusFailed
	default:
		dto.Status = dtoStatusUnknown
	}

	return dto
}

func JobsToDTO(jobs []entities.EvaluationJob) []JobDTO {
	result := make([]JobDTO, 0, len(jobs))

	for _, job := range jobs {
		result = append(result, JobToDTO(job))
	}

	return result
}
