package router

import "github.com/labstack/echo/v4"

type Router interface {
	AddRoutes(g *echo.Group)
}
