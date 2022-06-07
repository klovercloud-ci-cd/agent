package dependency

import (
	"github.com/klovercloud-ci-cd/agent/config"
	"github.com/klovercloud-ci-cd/agent/core/v1/logic"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

// GetV1ResourceService returns Resource service
func GetV1ResourceService() service.Resource {
	var observers []service.Observer
	eventStoreLogEventService := logic.NewEventStoreLogEventService(logic.NewHttpClientService())
	eventStoreProcessEventService := logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
	eventStoreProcessLifeCycleEvent := logic.NewEventStoreProcessLifeCycleService(logic.NewHttpClientService())
	observers = append(observers, eventStoreLogEventService)
	observers = append(observers, eventStoreProcessEventService)
	observers = append(observers, eventStoreProcessLifeCycleEvent)
	k8sClientSet, dynamicClient, discoveryClient := config.GetClientSet()
	var kubeEventPublisher service.KubeEventPublisher
	if config.KafkaPublisherEnabled {
		kubeEventPublisher = logic.NewKubeEventKafkaPublisher()
	} else {
		kubeEventPublisher = logic.NewKubeEventHttpPublisher(logic.NewHttpClientService())
	}
	return logic.NewResourceService(logic.NewK8sService(k8sClientSet, dynamicClient, discoveryClient, observers, kubeEventPublisher), observers, logic.NewHttpClientService())
}

// GetV1JwtService returns Jwt services
func GetV1JwtService() service.Jwt {
	return logic.NewJwtService()
}

// GetV1KubeEventService returns KubeEvent service
func GetV1KubeEventService() service.KubeEvent {
	var observers []service.Observer
	eventStoreLogEventService := logic.NewEventStoreLogEventService(logic.NewHttpClientService())
	eventStoreProcessEventService := logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
	eventStoreProcessLifeCycleEvent := logic.NewEventStoreProcessLifeCycleService(logic.NewHttpClientService())
	observers = append(observers, eventStoreLogEventService)
	observers = append(observers, eventStoreProcessEventService)
	observers = append(observers, eventStoreProcessLifeCycleEvent)
	k8sClientSet, dynamicClient, discoveryClient := config.GetClientSet()
	var kubeEventPublisher service.KubeEventPublisher
	if config.KafkaPublisherEnabled {
		kubeEventPublisher = logic.NewKubeEventKafkaPublisher()
	} else {
		kubeEventPublisher = logic.NewKubeEventHttpPublisher(logic.NewHttpClientService())
	}
	return logic.NewKubeEventService(logic.NewK8sService(k8sClientSet, dynamicClient, discoveryClient, observers, kubeEventPublisher))
}



// GetV1HttpClient returns HttpClient service
func GetV1HttpClient()service.HttpClient{
	return logic.NewHttpClientService()
}