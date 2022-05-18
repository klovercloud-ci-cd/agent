package main

import (
	"github.com/klovercloud-ci-cd/agent/api"
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/dependency"
	_ "github.com/klovercloud-ci-cd/agent/docs"
	"github.com/labstack/echo-contrib/jaegertracing"
	"io"
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
	api.Routes(e)
	resourceService := dependency.GetV1ResourceService()
	go ContinuePullingAgentEvent(resourceService)
	kubeEventService := dependency.GetV1KubeEventService()
	go kubeEventService.GetK8sObjectChangeEvents()
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

// ContinuePullingAgentEvent routine that pulls deploy steps in every interval.
func ContinuePullingAgentEvent(service service.Resource) {
	service.Pull()
	time.Sleep(time.Second * 5)
	ContinuePullingAgentEvent(service)
}
