package main

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/api"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/dependency"
	_ "github.com/klovercloud-ci-cd/agent/docs"
	"github.com/labstack/echo-contrib/jaegertracing"
	"io"
	"log"
	"time"
)

// @title agent API
// @description agent API
func main() {
	e := config.New()
	if config.EnableOpenTracing {
		c := jaegertracing.New(e, nil)
		defer func(c io.Closer) {
			err := c.Close()
			if err != nil {
				panic(err)
			}
		}(c)
	}
	go JoinIntegrationManager(dependency.GetV1HttpClient())
	api.Routes(e)
	resourceService := dependency.GetV1ResourceService()
	go ContinuePullingAgentEvent(resourceService)
	if config.LighthouseEnabled {
		kubeEventService := dependency.GetV1KubeEventService()
		go kubeEventService.GetK8sObjectChangeEvents()
	}
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

// ContinuePullingAgentEvent routine that pulls deploy steps in every interval.
func ContinuePullingAgentEvent(service service.Resource) {
	service.Pull()
	time.Sleep(time.Millisecond*500)
	ContinuePullingAgentEvent(service)
}

// JoinIntegrationManager joins clusters terminal with integration manager
func JoinIntegrationManager(client service.HttpClient) error {
	data := v1.Agent{
		Name:            config.AgentName,
		ApiVersion:      config.TerminalApiVersion,
		TerminalBaseUrl: config.TerminalBaseUrl,
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["token"] = config.Token
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = client.Post(config.ApiServiceUrl+"/agents", header, bytes)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
