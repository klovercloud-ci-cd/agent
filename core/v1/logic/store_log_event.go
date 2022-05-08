package logic

import (
	"encoding/json"
	"fmt"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"log"
	"strconv"
	"time"
)

type eventStoreEventService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreEventService) Listen(subject v1.Subject) {
	claim,_:=strconv.Atoi(fmt.Sprint(subject.EventData["claim"]))
	data := v1.LogEvent{
		ProcessId: subject.ProcessId,
		Log:       subject.Log,
		Step:      subject.Step,
		Footmark:  fmt.Sprint(subject.EventData["footmark"]),
		CreatedAt: time.Time{}.UTC(),
		Claim: claim,
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + config.Token
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
	e.httpPublisher.Post(config.ApiServiceUrl+"/logs", header, b)
}

// NewEventStoreLogEventService returns Observer type service
func NewEventStoreLogEventService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreEventService{
		httpPublisher: httpPublisher,
	}
}
