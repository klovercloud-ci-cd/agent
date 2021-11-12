package main

import (
	"github.com/klovercloud-ci-cd/agent/api"
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/dependency"
	_ "github.com/klovercloud-ci-cd/agent/docs"
	"time"
)

// @title agent API
// @description agent API
func main() {
	e := config.New()
	api.Routes(e)
	resourceService := dependency.GetV1ResourceService()
	go ContinuePullingAgentEvent(resourceService)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

// ContinuePullingAgentEvent routine that pulls deploy steps in every interval.
func ContinuePullingAgentEvent(service service.Resource) {
	service.Pull()
	time.Sleep(time.Second * 5)
	ContinuePullingAgentEvent(service)
}
