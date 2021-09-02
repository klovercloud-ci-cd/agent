package api

import (
	"github.com/labstack/echo/v4"
)

type Resource interface {
	Update(ctx echo.Context) error
}
