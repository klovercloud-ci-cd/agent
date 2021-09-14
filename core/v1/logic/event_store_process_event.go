package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/core/v1"
	"log"
)

type eventStoreProcessService struct {
	httpPublisher service.HttpPublisher
}

func (e eventStoreProcessService) Listen(subject v1.Subject) {
	if subject.EventData!=nil{
		event:=v1.PipelineProcessEvent{
			ProcessId: subject.ProcessId,
			Data:      subject.EventData,
		}
		header:=make(map[string]string)
		header["Content-Type"]="application/json"
		b, err := json.Marshal(event)
		if err!=nil{
			log.Println(err.Error())
			return
		}
		e.httpPublisher.Post(config.EventStoreUrl+"/processes",header,b)
	}
}

func NewEventStoreProcessEventService(httpPublisher service.HttpPublisher) service.EventStoreProcessEvent {
	return &eventStoreProcessService{
		httpPublisher: httpPublisher,
	}
}
