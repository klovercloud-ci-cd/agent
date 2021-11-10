package service

import v1 "github.com/klovercloud-ci-cd/agent/core/v1"

type Resource interface {
	Update(resource v1.Resource) error
	Pull()
}
