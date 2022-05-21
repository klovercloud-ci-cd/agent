package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kubeEventHttpPublisher struct {
	httpPublisher      service.HttpClient
}

func (k kubeEventHttpPublisher) Publish(message v1.KubeEventMessage) {
	marshal, _ := json.Marshal(message)
	header := make(map[string]string)
	header["token"] = config.Token
	header["Content-Type"] = "application/json"
	_ = k.httpPublisher.Post(config.ApiServiceUrl+"/kube_events", header, marshal)
	return
}

func NewKubeEventHttpPublisher(httpPublisher service.HttpClient) service.KubeEventPublisher {
	return kubeEventHttpPublisher{
		httpPublisher: httpPublisher,
	}
}

