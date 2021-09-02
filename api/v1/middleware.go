package v1

import (
	"github.com/labstack/echo/v4"
	"log"
)

// handle user authentication and authorization here.
func AuthenticationAndAuthorizationHandler(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
	log.Println("No authentication and authorization required ... ")
		return handler(context)
	}
}
