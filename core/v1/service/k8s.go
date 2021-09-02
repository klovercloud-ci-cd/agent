package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

type K8s interface {
	GetDeployment(name,namespace string)(*apiV1.Deployment,error)
	GetPod(name,namespace string)(*coreV1.Pod,error)
	GetStatefulSet(name,namespace string)(*apiV1.StatefulSet,error)
	GetDaemonSet(name,namespace string)(*apiV1.DaemonSet,error)
	UpdateDeployment(resource v1.Resource)error
	UpdatePod(resource v1.Resource)error
	UpdateStatefulSet(resource v1.Resource)error
	UpdateDaemonSet(resource v1.Resource)error
}

