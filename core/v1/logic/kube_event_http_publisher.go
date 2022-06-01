package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/enums"
	"log"
)

type kubeEventHttpPublisher struct {
	httpPublisher service.HttpClient
}

func (k kubeEventHttpPublisher) Publish(message v1.KubeEventMessage) {
	marshal, _ := json.Marshal(message)
	if message.Header.Command==enums.UPDATE{
		if message.Header.Extras["object"]=="pod"{
			log.Println(string(marshal))
		}
	}
	header := make(map[string]string)
	header["token"] = config.Token
	header["Content-Type"] = "application/json"
	err := k.httpPublisher.Post(config.ApiServiceUrl+"/kube_events", header, marshal)
	if err!=nil {
		log.Println(err.Error())
	}
	return
}

func NewKubeEventHttpPublisher(httpPublisher service.HttpClient) service.KubeEventPublisher {
	return kubeEventHttpPublisher{
		httpPublisher: httpPublisher,
	}
}
