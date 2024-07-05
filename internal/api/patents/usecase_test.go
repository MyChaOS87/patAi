//nolint:funlen // Test functions are long, due to test cases
package patents_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MyChaOS87/patAi/internal/api/patents"
	"github.com/MyChaOS87/patAi/internal/api/patents/mocks"
	"github.com/MyChaOS87/patAi/internal/authorization"
	"github.com/MyChaOS87/patAi/internal/entities"
)

var errFoo = errors.New("foo error")

type identity struct {
	id string
}

func (i *identity) GetID() string {
	return i.id
}

func Test_valuationJobUseCase_GetPatentValuationJobsByIdentityAndID(t *testing.T) {
	t.Parallel()

	var (
		alicesJob = entities.EvaluationJob{
			ID:                  uuid.MustParse("0441f94b-9a04-4015-9190-f213d55bf9fb"),
			OwnerID:             "Alice",
			EvaluationJobStatus: entities.EvaluationJobStatusPending,
			Value:               0,
		}
		bobsJob = entities.EvaluationJob{
			ID:                  uuid.MustParse("c32c6f83-b06b-4df5-9def-36e8ee5e6cb7"),
			OwnerID:             "Bob",
			EvaluationJobStatus: entities.EvaluationJobStatusFinished,
			Value:               42,
		}
	)

	testCases := []struct {
		name            string
		mockExpectation func(*mocks.QueueService)
		identity        authorization.Identity
		id              uuid.UUID
		want            entities.EvaluationJob
		wantErr         error
	}{
		{
			name: "Alice gets her job",
			mockExpectation: func(m *mocks.QueueService) {
				m.On("GetJobByID", uuid.MustParse("0441f94b-9a04-4015-9190-f213d55bf9fb")).Return(alicesJob, nil).Once()
			},
			identity: &identity{
				id: "Alice",
			},
			id: uuid.MustParse("0441f94b-9a04-4015-9190-f213d55bf9fb"),
			want: entities.EvaluationJob{
				ID:                  uuid.MustParse("0441f94b-9a04-4015-9190-f213d55bf9fb"),
				OwnerID:             "Alice",
				EvaluationJobStatus: entities.EvaluationJobStatusPending,
				Value:               0,
			},
			wantErr: nil,
		},
		{
			name: "Alice does not get Bob's job",
			mockExpectation: func(m *mocks.QueueService) {
				m.On("GetJobByID", uuid.MustParse("c32c6f83-b06b-4df5-9def-36e8ee5e6cb7")).Return(bobsJob, nil).Once()
			},
			identity: &identity{
				id: "Alice",
			},
			id:      uuid.MustParse("c32c6f83-b06b-4df5-9def-36e8ee5e6cb7"),
			want:    entities.EvaluationJob{},
			wantErr: patents.ErrJobNotFound,
		},
		{
			name: "Alice does not get a non-existent job",
			mockExpectation: func(m *mocks.QueueService) {
				m.On("GetJobByID", mock.Anything).Return(entities.EvaluationJob{}, patents.ErrJobNotFound).Once()
			},
			identity: &identity{
				id: "Alice",
			},
			id:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			want:    entities.EvaluationJob{},
			wantErr: patents.ErrJobNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			queueService := new(mocks.QueueService)

			tc.mockExpectation(queueService)

			useCase := patents.NewValuationJobUseCase(queueService, nil)

			job, err := useCase.GetPatentValuationJobByIdentityAndID(tc.identity, tc.id)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, job)

			queueService.AssertExpectations(t)
		})
	}
}

func Test_valuationJobUseCase_CreatePatentValuationJob(t *testing.T) {
	t.Parallel()

	content := "This is a patent content"

	alicesJob := entities.EvaluationJob{
		ID:                  uuid.MustParse("0441f94b-9a04-4015-9190-f213d55bf9fb"),
		OwnerID:             "Alice",
		EvaluationJobStatus: entities.EvaluationJobStatusPending,
		PatentContent:       content,
		Value:               0,
	}

	testCases := []struct {
		name        string
		preparation func(*mocks.QueueService, *mocks.QuotaService)
		identity    authorization.Identity
		want        entities.EvaluationJob
		wantErr     error
	}{
		{
			name: "Alice creates a job",
			preparation: func(queueService *mocks.QueueService, quotaService *mocks.QuotaService) {
				queueService.On("EnqueueJob", "Alice", content).Return(alicesJob, nil).Once()
				quotaService.On("GetQuotaToken", "Alice").Return(uuid.New(), nil).Once()
			},
			identity: &identity{
				id: "Alice",
			},
			want:    alicesJob,
			wantErr: nil,
		},
		{
			name: "Alice is over quota",
			preparation: func(_ *mocks.QueueService, quotaService *mocks.QuotaService) {
				quotaService.On("GetQuotaToken", "Alice").Return(uuid.Nil, patents.ErrQuotaExceeded).Once()
			},
			identity: &identity{
				id: "Alice",
			},
			want:    entities.EvaluationJob{},
			wantErr: patents.ErrQuotaExceeded,
		},
		{
			name: "Quota is returned if enqueue fails",
			preparation: func(queueService *mocks.QueueService, quotaService *mocks.QuotaService) {
				uuid := uuid.MustParse("e9f4ae48-a8bb-4c86-8530-f5756143480e")
				queueService.On("EnqueueJob", "Alice", content).Return(entities.EvaluationJob{}, errFoo).Once()
				quotaService.On("GetQuotaToken", "Alice").Return(uuid, nil).Once()
				quotaService.On("ReturnQuotaToken", uuid).Return().Once()
			},
			identity: &identity{
				id: "Alice",
			},
			want:    entities.EvaluationJob{},
			wantErr: errFoo,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			queueService := new(mocks.QueueService)
			quotaService := new(mocks.QuotaService)

			tc.preparation(queueService, quotaService)

			useCase := patents.NewValuationJobUseCase(queueService, quotaService)

			job, err := useCase.CreatePatentValuationJob(tc.identity, content)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, job)

			queueService.AssertExpectations(t)
			quotaService.AssertExpectations(t)
		})
	}
}
