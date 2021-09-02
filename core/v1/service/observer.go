package service
import (
	v1 "github.com/klovercloud-ci/core/v1"
)
type Observer interface {
	Listen(v1.Subject)
}
