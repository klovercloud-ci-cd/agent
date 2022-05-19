package logic

import (
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kubeEventHttpPublisher struct {
}

func (k kubeEventHttpPublisher) Publish(message v1.KubeEventMessage) {
	panic("implement me")
}

func NewKubeEventHttpPublisher() service.KubeEventPublisher {
	return kubeEventHttpPublisher{}
}

