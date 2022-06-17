package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/enums"
	"log"
	"time"
)

type eventStoreProcessLifeCycleService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreProcessLifeCycleService) Listen(subject v1.Subject) {
	if subject.EventData == nil {
		return
	}
	if subject.EventData != nil {
		data := []v1.ProcessLifeCycleEvent{}
		processLifeCycleEvent := v1.ProcessLifeCycleEvent{
			ProcessId: subject.ProcessId,
			Step:      subject.Step,
			Agent:     config.AgentName,
			CreatedAt: time.Now().UTC(),
		}
		nextSteps := []string{}

		if subject.Pipeline != nil {
			for _, step := range subject.Pipeline.Steps {
				if step.Name == subject.Step {
					nextSteps = append(nextSteps, step.Next...)
				}
			}
		}
		if subject.EventData["status"] == enums.DEPLOYMENT_FAILED || subject.EventData["status"] == enums.ERROR || subject.EventData["status"] == enums.TERMINATING || subject.EventData["status"] ==enums.FAILED {
			processLifeCycleEvent.Status = enums.FAILED
			data = append(data, processLifeCycleEvent)
		} else if subject.EventData["status"] == enums.SUCCESSFUL || subject.EventData["status"] == enums.COMPLETED{
			log.Println(subject.EventData["status"])
			processLifeCycleEvent.Status = enums.COMPLETED
			data = append(data, processLifeCycleEvent)
			for _, each := range nextSteps {
				data = append(data, v1.ProcessLifeCycleEvent{
					ProcessId: subject.ProcessId,
					Status:    enums.PAUSED,
					Step:      each,
					Agent:     config.AgentName,
					CreatedAt: time.Now().UTC(),
				})
			}
		}else{
			return
		}
		type ProcessLifeCycleEventList struct {
			Events []v1.ProcessLifeCycleEvent `bson:"events" json:"events"`
		}
		if len(data) > 0 {
			events := ProcessLifeCycleEventList{data}
			header := make(map[string]string)
			header["Content-Type"] = "application/json"
			header["token"] = config.Token
			b, err := json.Marshal(events)
			if err != nil {
				log.Println(err.Error())
				return
			}
			e.httpPublisher.Post(config.ApiServiceUrl+"/process_life_cycle_events", header, b)
		}
	}
}

// NewEventStoreProcessLifeCycleService returns Observer type service
func NewEventStoreProcessLifeCycleService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreProcessLifeCycleService{
		httpPublisher: httpPublisher,
	}
}
