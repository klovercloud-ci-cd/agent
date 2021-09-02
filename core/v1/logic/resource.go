package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
)

type resourceService struct {
	K8s          service.K8s
	observerList []service.Observer
}

func (r resourceService) Update(resource v1.Resource) error {
	if resource.Type == enums.DEPLOYMENT {
		return r.K8s.UpdateDeployment(resource)
	} else if resource.Type == enums.POD {
		return r.K8s.UpdatePod(resource)
	} else if resource.Type == enums.STATEFULSET {
		return r.K8s.UpdateStatefulSet(resource)
	} else if resource.Type == enums.DAEMONSET {
		return r.K8s.UpdateDaemonSet(resource)
	}
	return nil
}
func (r resourceService)notifyAll(subject v1.Subject){
	for _, observer := range r.observerList {
		go observer.Listen(subject)
	}
}

func NewResourceService(k8s service.K8s, observerList []service.Observer) service.Resource {
	return &resourceService{
		K8s:          k8s,
		observerList: observerList,
	}
}
