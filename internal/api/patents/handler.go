package patents

import (
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/internal/authorization"
	"github.com/MyChaOS87/patAi/pkg/log"
)

type handler struct {
	useCase ValuationJobUseCase
}

func NewHandler(useCase ValuationJobUseCase) Handler {
	return &handler{
		useCase: useCase,
	}
}

var errGetIdentityFailed = errors.New("cannot get identity from context")

func getIdentityFromContext(c echo.Context) (authorization.Identity, error) {
	identity, ok := c.Get(contextIdentityKey).(authorization.Identity)
	if !ok {
		return nil, errGetIdentityFailed
	}

	return identity, nil
}

func (h *handler) GetPatentValuationJobs() echo.HandlerFunc {
	return func(c echo.Context) error {
		identity, err := getIdentityFromContext(c)
		if err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		jobs, err := h.useCase.GetPatentValuationJobsByIdentity(identity)
		if err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := c.JSON(http.StatusOK, JobsToDTO(jobs)); err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (h *handler) GetPatentValuationJobByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		identity, err := getIdentityFromContext(c)
		if err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		uuid, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "malformed job id")
		}

		job, err := h.useCase.GetPatentValuationJobByIdentityAndID(identity, uuid)
		if errors.Is(err, ErrJobNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "job not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := c.JSON(http.StatusOK, JobToDTO(job)); err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (h *handler) CreatePatentValuationJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		identity, err := getIdentityFromContext(c)
		if err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		body := new(strings.Builder)
		if _, err := io.Copy(body, c.Request().Body); err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		content := body.String()

		job, err := h.useCase.CreatePatentValuationJob(identity, content)
		if errors.Is(err, ErrQuotaExceeded) {
			return echo.NewHTTPError(http.StatusTooManyRequests, err.Error())
		} else if err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := c.JSON(http.StatusCreated, JobToDTO(job)); err != nil {
			log.Errorf("%v", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}
