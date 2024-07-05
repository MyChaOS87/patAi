package simulation

import (
	"time"

	"github.com/google/uuid"

	"github.com/MyChaOS87/patAi/internal/api/patents"
	"github.com/MyChaOS87/patAi/internal/entities"
	"github.com/MyChaOS87/patAi/pkg/log"
)

type Simulation interface {
	patents.QueueService
	patents.QuotaService
}

type inMemoryQueueAndQuotaServiceSimulation struct {
	jobs               []*entities.EvaluationJob
	jobsByID           map[uuid.UUID]*entities.EvaluationJob
	jobsByOwner        map[string][]*entities.EvaluationJob
	quotaTokensByOwner map[string][]uuid.UUID
}

func NewInMemoryQueueAndQuotaServiceSimulation() Simulation {
	return &inMemoryQueueAndQuotaServiceSimulation{
		jobs:               []*entities.EvaluationJob{},
		jobsByID:           map[uuid.UUID]*entities.EvaluationJob{},
		jobsByOwner:        map[string][]*entities.EvaluationJob{},
		quotaTokensByOwner: map[string][]uuid.UUID{},
	}
}

func (s *inMemoryQueueAndQuotaServiceSimulation) EnqueueJob(ownerID string, content string) (entities.EvaluationJob, error) {
	job := entities.EvaluationJob{
		ID:                  uuid.New(),
		OwnerID:             ownerID,
		EvaluationJobStatus: entities.EvaluationJobStatusPending,
		PatentContent:       content,
	}

	s.jobs = append(s.jobs, &job)
	s.jobsByID[job.ID] = &job
	s.jobsByOwner[job.OwnerID] = append(s.jobsByOwner[job.OwnerID], &job)

	log.Infof("Job %s scheduled for execution", job.ID.String())

	// Simulate trigger of evaluation
	go func() {
		//nolint:gomnd // Simulation preset
		time.Sleep(2 * time.Minute)

		log.Infof("Job %s finished evaluation", job.ID.String())
		job.EvaluationJobStatus = entities.EvaluationJobStatusFinished
		job.Value = 42
	}()

	return job, nil
}

func (s *inMemoryQueueAndQuotaServiceSimulation) GetJobsByOwnerID(ownerID string) ([]entities.EvaluationJob, error) {
	jobs := s.jobsByOwner[ownerID]
	if jobs == nil {
		return nil, nil
	}

	result := make([]entities.EvaluationJob, len(jobs))
	for i, j := range jobs {
		result[i] = *j
	}

	return result, nil
}

func (s *inMemoryQueueAndQuotaServiceSimulation) GetJobByID(id uuid.UUID) (entities.EvaluationJob, error) {
	job := s.jobsByID[id]
	if job == nil {
		return entities.EvaluationJob{}, patents.ErrJobNotFound
	}

	return *job, nil
}

func (s *inMemoryQueueAndQuotaServiceSimulation) GetQuotaToken(ownerID string) (uuid.UUID, error) {
	tokens := s.quotaTokensByOwner[ownerID]
	//nolint:gomnd // Simulation preset
	if len(tokens) >= 5 {
		return uuid.UUID{}, patents.ErrQuotaExceeded
	}

	token := uuid.New()
	s.quotaTokensByOwner[ownerID] = append(tokens, token)

	// Simulate token expiration after 5 min
	go func() {
		//nolint:gomnd // Simulation preset
		time.Sleep(5 * time.Minute)

		s.ReturnQuotaToken(token)
	}()

	return token, nil
}

func (s *inMemoryQueueAndQuotaServiceSimulation) ReturnQuotaToken(token uuid.UUID) {
	for ownerID, tokens := range s.quotaTokensByOwner {
		for i, t := range tokens {
			if t == token {
				s.quotaTokensByOwner[ownerID] = append(tokens[:i], tokens[i+1:]...)

				return
			}
		}
	}
}
