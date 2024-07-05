//go:generate mockery --name QueueService|QuotaService

package patents

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/internal/entities"
)

var (
	ErrQuotaExceeded         = errors.New("quota exceeded")
	ErrCouldNotRetrieveQuota = errors.New("could not retrieve quota token")
	ErrCouldNotEnqueueJob    = errors.New("could not enqueue job")
	ErrJobNotFound           = errors.New("job not found")
)

type QueueService interface {
	EnqueueJob(ownerID string, content string) (entities.EvaluationJob, error)
	GetJobsByOwnerID(ownerID string) ([]entities.EvaluationJob, error)
	GetJobByID(id uuid.UUID) (entities.EvaluationJob, error)
}

type QuotaService interface {
	// returns a token that can be used to enqueue a job and returns an ErrQuotaExceeded error if the quota is exceeded
	GetQuotaToken(ownerID string) (uuid.UUID, error)
	ReturnQuotaToken(token uuid.UUID)
}
