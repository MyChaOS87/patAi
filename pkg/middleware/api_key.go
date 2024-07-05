package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/pkg/log"
)

type AuthorizationProvider[T interface{}] interface {
	GetByAPIKey(apiKey string) (T, error)
}

//nolint:gosec // false positive no credentials here;
const apiKeyHeaderField = "X-API-Key"

func mapAPIKeyError400To401(inner echo.MiddlewareFunc) echo.MiddlewareFunc {
	errorToReplace := echo.NewHTTPError(http.StatusBadRequest, "missing key in request header")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := inner(next)(c)

			if err == nil {
				return nil
			}

			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) && httpErr.Code == errorToReplace.Code && httpErr.Message == errorToReplace.Message {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing key X-API-Key in request header")
			}

			return err
		}
	}
}

func APIKey[T interface{}](authorizationProvider AuthorizationProvider[T],
	contextIdentityKey string,
) echo.MiddlewareFunc {
	const errMessage = "Unauthorized: API-Key auth failed"

	return mapAPIKeyError400To401(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: fmt.Sprintf("header:%s", apiKeyHeaderField),
		Validator: func(apiKey string, c echo.Context) (bool, error) {
			identity, err := authorizationProvider.GetByAPIKey(apiKey)
			if err != nil {
				log.Errorf("api key identity lookup failed: %v", err)

				return false, errors.Wrap(err, errMessage)
			}

			c.Set(contextIdentityKey, identity)

			return true, nil
		},
	}))
}
