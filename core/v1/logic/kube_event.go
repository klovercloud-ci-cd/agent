package logic

import (
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
)

type kubeEventService struct {
	k8s service.K8s
}

func (k kubeEventService) GetK8sObjectChangeEvents() {
	_, namespace_controller := k.k8s.ListenNamespaceEvents()
	_, pod_controller := k.k8s.ListenPodEvents()
	_, deployment_controller := k.k8s.ListenDeployEvents()
	_, service_controller := k.k8s.ListenServiceEvents()
	_, replicaSet_controller := k.k8s.ListenReplicaSetEvents()
	_, stateFullSet_controller := k.k8s.ListenStateFullSetSetEvents()
	_, daemonSet_controller := k.k8s.ListenDaemonSetEvents()
	_, pv_controller := k.k8s.ListenPVEvents()
	_, pvc_controller := k.k8s.ListenPVCEvents()
	_, cm_controller := k.k8s.ListenConfigMapEvents()
	_, secret_controller := k.k8s.ListenSecretEvents()
	_, serviceAccount_controller := k.k8s.ListenServiceAccountEvents()
	_, roles_controller := k.k8s.ListenRoleEvents()
	_, rolebindings_controller := k.k8s.ListenRoleBindingEvents()
	_, clusterrole_controller := k.k8s.ListenClusterRoleEvents()
	_, clusterrolebindings_controller := k.k8s.ListenClusterRoleBindingEvents()
	_, networkPolicy_controller := k.k8s.ListenNetworkPolicyEvents()
	_, ingress_controller := k.k8s.ListenIngressEvents()

	stop := make(chan struct{})

	go namespace_controller.Run(stop)
	go pod_controller.Run(stop)
	go deployment_controller.Run(stop)
	go service_controller.Run(stop)
	go replicaSet_controller.Run(stop)
	go stateFullSet_controller.Run(stop)
	go daemonSet_controller.Run(stop)
	go pv_controller.Run(stop)
	go pvc_controller.Run(stop)
	go cm_controller.Run(stop)
	go secret_controller.Run(stop)
	go serviceAccount_controller.Run(stop)
	go roles_controller.Run(stop)
	go rolebindings_controller.Run(stop)
	go clusterrole_controller.Run(stop)
	go clusterrolebindings_controller.Run(stop)
	go networkPolicy_controller.Run(stop)
	go ingress_controller.Run(stop)
}

func NewKubeEventService(k8s service.K8s) service.KubeEvent {
	return &kubeEventService{
		k8s: k8s,
	}
}
