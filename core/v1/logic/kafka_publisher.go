package logic

import (
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kafkaPublisher struct {
}

func (k kafkaPublisher) Publish(message v1.KafkaMessage) {
	panic("implement me")
}

func NewKafkaPublisher() service.KafkaPublisher {
	return kafkaPublisher{}
}
