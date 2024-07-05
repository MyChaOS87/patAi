package patents

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/internal/authorization"
	"github.com/MyChaOS87/patAi/internal/entities"
)

var ErrValuationUseCase = errors.New("valuation use case error")

type ValuationJobUseCase interface {
	GetPatentValuationJobsByIdentity(identity authorization.Identity) ([]entities.EvaluationJob, error)
	GetPatentValuationJobByIdentityAndID(identity authorization.Identity, ID uuid.UUID) (entities.EvaluationJob, error)

	// CreatePatentValuationJob creates a new patent valuation job after checking the users quota
	// returns an ErrQuotaExceeded error if the user has exceeded their quota
	CreatePatentValuationJob(identity authorization.Identity, content string) (entities.EvaluationJob, error)
}

type valuationJobUseCase struct {
	queueService QueueService
	quotaService QuotaService
}

func NewValuationJobUseCase(queueService QueueService, quotaService QuotaService) ValuationJobUseCase {
	return &valuationJobUseCase{
		queueService: queueService,
		quotaService: quotaService,
	}
}

func (v *valuationJobUseCase) GetPatentValuationJobsByIdentity(
	identity authorization.Identity,
) ([]entities.EvaluationJob, error) {
	res, err := v.queueService.GetJobsByOwnerID(identity.GetID())
	if err != nil {
		return nil, errors.Wrap(err, ErrValuationUseCase.Error())
	}

	return res, nil
}

func (v *valuationJobUseCase) GetPatentValuationJobByIdentityAndID(
	identity authorization.Identity,
	id uuid.UUID,
) (entities.EvaluationJob, error) {
	job, err := v.queueService.GetJobByID(id)
	if err != nil {
		return entities.EvaluationJob{}, errors.Wrap(err, ErrValuationUseCase.Error())
	}

	if job.OwnerID != identity.GetID() {
		return entities.EvaluationJob{}, ErrJobNotFound
	}

	return job, nil
}

func (v *valuationJobUseCase) CreatePatentValuationJob(
	identity authorization.Identity, content string,
) (entities.EvaluationJob, error) {
	token, err := v.quotaService.GetQuotaToken(identity.GetID())
	if err != nil {
		return entities.EvaluationJob{}, errors.Wrap(err, ErrCouldNotRetrieveQuota.Error())
	}

	job, err := v.queueService.EnqueueJob(identity.GetID(), content)
	if err != nil {
		v.quotaService.ReturnQuotaToken(token)

		return entities.EvaluationJob{}, errors.Wrap(err, ErrValuationUseCase.Error())
	}

	return job, nil
}
