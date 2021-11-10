package service

import (
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
)

// Observer Observer operations.
type Observer interface {
	Listen(v1.Subject)
}
