package service

import v1 "github.com/klovercloud-ci/core/v1"

type Resource interface {
	Update(resource v1.Resource) error
}
