package v1

import (
	"github.com/klovercloud-ci/dependency"
	"github.com/labstack/echo/v4"
)

func Router(g *echo.Group) {
	ResourceRouter(g.Group("/resources"))
}

func ResourceRouter(g *echo.Group) {
	resourceRouter := NewResourceApi(dependency.GetResourceService())
	g.POST("", resourceRouter.Update, AuthenticationAndAuthorizationHandler)
}
