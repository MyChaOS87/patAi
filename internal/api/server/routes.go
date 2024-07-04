package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/pkg/openapi"
)

const (
	v0BaseURI = "/api/v0/"
	v0Health  = "health"
	V0OpenAPI = "openapi"

	bodyLimit = "2M"
)

func (s *Server) mapHandlers() error {
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit(bodyLimit))
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.api.AllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}))

	v0 := s.echo.Group(v0BaseURI)

	// health
	s.mapHealthRoute(v0.Group(v0Health))

	// Documentation
	openapi.MapDocumentationRoutes(
		v0.Group(V0OpenAPI),
		openapi.NewHandlers(),
		s.api.OpenAPIFile, s.api.OpenAPISwaggerUI, struct{ ServerBaseURL string }{ServerBaseURL: s.api.ServerBaseURL})

	for _, r := range s.childRouters {
		r.AddRoutes(v0)
	}

	return nil
}

func (s *Server) mapHealthRoute(g *echo.Group) {
	g.GET("", func(c echo.Context) error {
		if err := c.JSON(http.StatusOK, map[string]string{"status": "OK"}); err != nil {
			return errors.Wrap(err, "health failed")
		}

		return nil
	})
}
