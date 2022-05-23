package logic

import (
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kubeEventService struct {
	k8s service.K8s
}

func (k kubeEventService) GetK8sObjectChangeEvents() {
	_, namespaceController := k.k8s.ListenNamespaceEvents()
	_, podController := k.k8s.ListenPodEvents()
	_, deploymentController := k.k8s.ListenDeployEvents()
	_, serviceController := k.k8s.ListenServiceEvents()
	_, replicasetController := k.k8s.ListenReplicaSetEvents()
	_, stateFullSetController := k.k8s.ListenStateFullSetSetEvents()
	_, daemonSetController := k.k8s.ListenDaemonSetEvents()
	_, pvController := k.k8s.ListenPVEvents()
	_, pvcController := k.k8s.ListenPVCEvents()
	_, cmController := k.k8s.ListenConfigMapEvents()
	_, secretController := k.k8s.ListenSecretEvents()
	_, serviceAccountController := k.k8s.ListenServiceAccountEvents()
	_, rolesController := k.k8s.ListenRoleEvents()
	_, roleBindingController := k.k8s.ListenRoleBindingEvents()
	_, clusterRoleController := k.k8s.ListenClusterRoleEvents()
	_, clusterRoleBindingController := k.k8s.ListenClusterRoleBindingEvents()
	_, networkPolicyController := k.k8s.ListenNetworkPolicyEvents()
	_, ingressController := k.k8s.ListenIngressEvents()
	_, kubeEventController := k.k8s.ListenKubeEvents()
	stop := make(chan struct{})
	go namespaceController.Run(stop)
	go podController.Run(stop)
	go deploymentController.Run(stop)
	go serviceController.Run(stop)
	go replicasetController.Run(stop)
	go stateFullSetController.Run(stop)
	go daemonSetController.Run(stop)
	go pvController.Run(stop)
	go pvcController.Run(stop)
	go cmController.Run(stop)
	go secretController.Run(stop)
	go serviceAccountController.Run(stop)
	go rolesController.Run(stop)
	go roleBindingController.Run(stop)
	go clusterRoleController.Run(stop)
	go clusterRoleBindingController.Run(stop)
	go networkPolicyController.Run(stop)
	go ingressController.Run(stop)
	go kubeEventController.Run(stop)
}

func NewKubeEventService(k8s service.K8s) service.KubeEvent {
	return &kubeEventService{
		k8s: k8s,
	}
}
