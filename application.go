package main

import (
	"github.com/klovercloud-ci/api"
	"github.com/klovercloud-ci/config"
)

func main(){
	e:=config.New()
	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}
