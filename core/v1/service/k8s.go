package service

import (
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

// K8s operations.
type K8s interface {
	GetDeployment(name, namespace string) (*apiV1.Deployment, error)
	GetPod(name, namespace string) (*coreV1.Pod, error)
	GetStatefulSet(name, namespace string) (*apiV1.StatefulSet, error)
	GetDaemonSet(name, namespace string) (*apiV1.DaemonSet, error)
	UpdateDeployment(resource v1.Resource) error
	UpdatePod(resource v1.Resource) error
	UpdateStatefulSet(resource v1.Resource) error
	UpdateDaemonSet(resource v1.Resource) error
	Apply(data unstructured.Unstructured) error
	Deploy(data *unstructured.Unstructured) (bool, error)
	ListenNamespaceEvents() (cache.Store, cache.Controller)
	ListenServiceEvents() (cache.Store, cache.Controller)
	ListenPodEvents() (cache.Store, cache.Controller)
	ListenDeployEvents() (cache.Store, cache.Controller)
	ListenIngressEvents() (cache.Store, cache.Controller)
	ListenNetworkPolicyEvents() (cache.Store, cache.Controller)
	ListenClusterRoleBindingEvents() (cache.Store, cache.Controller)
	ListenClusterRoleEvents() (cache.Store, cache.Controller)
	ListenRoleBindingEvents() (cache.Store, cache.Controller)
	ListenRoleEvents() (cache.Store, cache.Controller)
	ListenServiceAccountEvents() (cache.Store, cache.Controller)
	ListenSecretEvents() (cache.Store, cache.Controller)
	ListenConfigMapEvents() (cache.Store, cache.Controller)
	ListenPVCEvents() (cache.Store, cache.Controller)
	ListenPVEvents() (cache.Store, cache.Controller)
	ListenDaemonSetEvents() (cache.Store, cache.Controller)
	ListenReplicaSetEvents() (cache.Store, cache.Controller)
	ListenStateFullSetSetEvents() (cache.Store, cache.Controller)
}
