package main

import (
	"github.com/klovercloud-ci-cd/agent/api"
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/dependency"
	"time"
)

func main(){
	e:=config.New()
	api.Routes(e)
	resourceService:=dependency.GetResourceService()
	go ContinuePullingAgentEvent(resourceService)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}
func ContinuePullingAgentEvent(service service.Resource){
	service.Pull()
	time.Sleep(time.Second*5)
	ContinuePullingAgentEvent(service)
}