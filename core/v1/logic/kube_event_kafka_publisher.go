package logic

import (
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kubeEventKafkaPublisher struct {
}

func (k kubeEventKafkaPublisher) Publish(message v1.KubeEventMessage) {
	panic("implement me")
}

func NewKafkaPublisher() service.KubeEventPublisher {
	return kubeEventKafkaPublisher{}
}
