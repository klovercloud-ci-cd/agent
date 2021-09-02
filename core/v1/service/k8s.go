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
	UpdateDeployment(deployment apiV1.Deployment,resource v1.Resource)error
	UpdatePod(deployment coreV1.Pod,resource v1.Resource)error
	UpdateStatefulSet(statefulSet apiV1.StatefulSet,resource v1.Resource)error
	UpdateDaemonSet(daemonSet apiV1.DaemonSet,resource v1.Resource)error
}

