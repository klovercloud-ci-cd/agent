package dependency

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"

)

func GetResourceService() service.Resource{
	var observers [] service.Observer
	eventStoreLogEventService:=logic.NewEventStoreLogEventService(logic.NewHttpClientService())
	eventStoreProcessEventService:=logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
	observers= append(observers, eventStoreLogEventService)
	observers= append(observers, eventStoreProcessEventService)
	k8sClientSet,dynamicClient,discoveryClient:=config.GetClientSet()
	return logic.NewResourceService(logic.NewK8sService(k8sClientSet,dynamicClient,discoveryClient,observers),observers,logic.NewHttpClientService())
}