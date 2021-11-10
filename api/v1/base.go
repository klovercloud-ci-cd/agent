package v1

import (
	"github.com/klovercloud-ci-cd/agent/dependency"
	"github.com/labstack/echo/v4"
)

// Router api/v1 base router
func Router(g *echo.Group) {
	ResourceRouter(g.Group("/resources"))
}

// ResourceRouter api/v1/resources/* router
func ResourceRouter(g *echo.Group) {
	resourceRouter := NewResourceApi(dependency.GetV1ResourceService())
	g.POST("", resourceRouter.Update, AuthenticationAndAuthorizationHandler)
}
