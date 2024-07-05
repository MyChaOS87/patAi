package patents

import (
	"github.com/labstack/echo/v4"

	"github.com/MyChaOS87/patAi/internal/api/router"
	"github.com/MyChaOS87/patAi/internal/authorization"
	"github.com/MyChaOS87/patAi/pkg/middleware"
)

const (
	patentsBaseURI     = "patents"
	contextIdentityKey = "patents-identity"
)

var _ router.Router = &patents{}

type Handler interface {
	GetPatentValuationJobs() echo.HandlerFunc
	GetPatentValuationJobByID() echo.HandlerFunc
	CreatePatentValuationJob() echo.HandlerFunc
}

type patents struct {
	authorizationProvider middleware.AuthorizationProvider[authorization.Identity]
	handler               Handler
}

func NewPatentsRouter(
	authorizationProvider middleware.AuthorizationProvider[authorization.Identity], handler Handler,
) router.Router {
	return &patents{
		authorizationProvider: authorizationProvider,
		handler:               handler,
	}
}

func (p *patents) AddRoutes(baseGroup *echo.Group) {
	patentsGroup := baseGroup.Group(patentsBaseURI)
	patentsGroup.Use(middleware.APIKey(p.authorizationProvider, contextIdentityKey))

	patentsGroup.GET("", p.handler.GetPatentValuationJobs())
	patentsGroup.GET("/:id", p.handler.GetPatentValuationJobByID())
	patentsGroup.POST("", p.handler.CreatePatentValuationJob())
}
