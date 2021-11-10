package api

import (
	"github.com/labstack/echo/v4"
)

// Resource resource api operations
type Resource interface {
	Update(ctx echo.Context) error
}
