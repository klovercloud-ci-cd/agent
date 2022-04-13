package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"log"
)

type eventStoreProcessService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreProcessService) Listen(subject v1.Subject) {
	if subject.EventData != nil {
		event := v1.PipelineProcessEvent{
			ProcessId: subject.ProcessId,
			Data:      subject.EventData,
		}
		header := make(map[string]string)
		header["Content-Type"] = "application/json"
		header["Authorization"] = "Bearer " + config.Token
		b, err := json.Marshal(event)
		if err != nil {
			log.Println(err.Error())
			return
		}
		e.httpPublisher.Post(config.ApiServiceUrl+"/processes_events", header, b)

	}
}

// NewEventStoreProcessEventService returns Observer type service
func NewEventStoreProcessEventService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreProcessService{
		httpPublisher: httpPublisher,
	}
}
