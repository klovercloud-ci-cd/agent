package main

import (
	"github.com/klovercloud-ci/api"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	"github.com/klovercloud-ci/core/v1/service"
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