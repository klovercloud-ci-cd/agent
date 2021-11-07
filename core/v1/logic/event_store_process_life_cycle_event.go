package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"log"
	"time"
)

type eventStoreProcessLifeCycleService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreProcessLifeCycleService) Listen(subject v1.Subject) {
	log.Println(subject.EventData["status"])
	if subject.EventData["status"]==nil{
		return
	}
	if subject.EventData != nil {
		data := []v1.ProcessLifeCycleEvent{}
		processLifeCycleEvent := v1.ProcessLifeCycleEvent{
			ProcessId: subject.ProcessId,
			Step:      subject.Step,
			CreatedAt: time.Now().UTC(),
		}
		nestStepNameMap := make(map[string]bool)
		for _, name := range processLifeCycleEvent.Next {
			nestStepNameMap[name] = true
		}
		 if subject.EventData["status"] == enums.DEPLOYMENT_FAILED || subject.EventData["status"] == enums.ERROR || subject.EventData["status"] == enums.TERMINATING{
			processLifeCycleEvent.Status = enums.FAILED
			data = append(data,processLifeCycleEvent)
		}else if subject.EventData["status"] == enums.SUCCESSFUL{
			processLifeCycleEvent.Status = enums.COMPLETED
			data = append(data,processLifeCycleEvent)
		}else{
			return
		 }
		type ProcessLifeCycleEventList struct {
			Events [] v1.ProcessLifeCycleEvent `bson:"events" json :"events"`
		}
		if len(data)>0 {
			events := ProcessLifeCycleEventList{data}
			header := make(map[string]string)
			header["Content-Type"] = "application/json"
			header["token"]=config.Token
			b, err := json.Marshal(events)
			if err != nil {
				log.Println(err.Error())
				return
			}
			e.httpPublisher.Post(config.EventStoreUrl+"/process_life_cycle_events", header, b)
		}
	}
}

func NewEventStoreProcessLifeCycleService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreProcessLifeCycleService{
		httpPublisher: httpPublisher,
	}
}
