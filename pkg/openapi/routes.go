package openapi

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers interface {
	GetOpenAPIDocument(openAPIFile string, template interface{}) echo.HandlerFunc
}

func MapDocumentationRoutes(
	g *echo.Group, h Handlers, openAPIFile string, serveSwaggerUI bool, template interface{},
) {
	g.GET("", h.GetOpenAPIDocument(openAPIFile, template))

	if serveSwaggerUI {
		g.GET("/openapi.yaml", h.GetOpenAPIDocument(openAPIFile, template))
		g.GET("/ui/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("../openapi.yaml")))
	}
}
