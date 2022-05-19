package service

import v1 "github.com/klovercloud-ci-cd/agent/core/v1"

type KubeEventPublisher interface {
	Publish(message v1.KubeEventMessage)
}
