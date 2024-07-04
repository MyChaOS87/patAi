package openapi

import (
	"bytes"
	"net/http"
	tmpl "text/template"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/pkg/log"
)

type HTTPHandlers struct{}

func NewHandlers() *HTTPHandlers {
	return &HTTPHandlers{}
}

func (h *HTTPHandlers) GetOpenAPIDocument(openAPIFile string, template interface{}) echo.HandlerFunc {
	openAPITemplate := tmpl.Must(tmpl.ParseGlob(openAPIFile))

	var openAPIRendered bytes.Buffer

	if err := openAPITemplate.Execute(&openAPIRendered, template); err != nil {
		log.Fatalf("cannot execute template: %v", err)
	}

	return func(c echo.Context) error {
		if err := c.Blob(http.StatusOK, "text/yaml", openAPIRendered.Bytes()); err != nil {
			log.Error(err)

			return errors.Wrap(err, "OpenAPI file failed")
		}

		return nil
	}
}
