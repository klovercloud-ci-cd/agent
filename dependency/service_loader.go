package dependency

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"

)

func GetResourceService() service.Resource{
	var observers [] service.Observer
	eventStoreLogEventService:=logic.NewEventStoreLogEventService(logic.NewHttpPublisherService())
	eventStoreProcessEventService:=logic.NewEventStoreProcessEventService(logic.NewHttpPublisherService())
	observers= append(observers, eventStoreLogEventService)
	observers= append(observers, eventStoreProcessEventService)
	k8sClientSet:=config.GetClientSet()
	return logic.NewResourceService(logic.NewK8sService(k8sClientSet,observers),observers)
}