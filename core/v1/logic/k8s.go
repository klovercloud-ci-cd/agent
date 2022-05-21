package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/enums"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	networkingV1 "k8s.io/api/networking/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"log"
	"strconv"
	"strings"
	"time"
)

type k8sService struct {
	kcs             *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	discoveryClient *discovery.DiscoveryClient
	observerList       []service.Observer
	kubeEventPublisher service.KubeEventPublisher
}

func (k k8sService) ListenNamespaceEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "namespaces", "",
		fields.Everything())
	extrasMap := map[string]string{"type":"namespace", "agent": config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.Namespace{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ns := obj.(*coreV1.Namespace)
				if _, ok := ns.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   ns,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Namespace: ", ns.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				ns := obj.(*coreV1.Namespace)
				if _, ok := ns.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   ns,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Namespace:", ns.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.Namespace)
				newK8sObj := newObj.(*coreV1.Namespace)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Namespace:", oldK8sObj.Name, ", new Namespace:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenServiceEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "services", "", fields.Everything())
	extrasMap := map[string]string{"type":"service", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				svc := obj.(*coreV1.Service)
				if _, ok := svc.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   svc,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Service: ", svc.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				svc := obj.(*coreV1.Service)
				if _, ok := svc.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   svc,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Service:", svc.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.Service)
				newK8sObj := newObj.(*coreV1.Service)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Service:", oldK8sObj.Name, ", new Service:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenPodEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "pods", "",
		fields.Everything())
	extrasMap := map[string]string{"type":"pod", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*coreV1.Pod)
				if _, ok := pod.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   pod,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Pod: ", pod.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*coreV1.Pod)
				if _, ok := pod.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   pod,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Pod:", pod.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.Pod)
				newK8sObj := newObj.(*coreV1.Pod)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Pod:", oldK8sObj.Name, ", new Pod:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenDeployEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.AppsV1().RESTClient(), "deployments", "", fields.Everything())
	extrasMap := map[string]string{"type":"deployment", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&appsV1.Deployment{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment := obj.(*appsV1.Deployment)
				if _, ok := deployment.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   deployment,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Deploy: ", deployment.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				deployment := obj.(*appsV1.Deployment)
				if _, ok := deployment.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   deployment,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Deploy: ", deployment.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*appsV1.Deployment)
				newK8sObj := newObj.(*appsV1.Deployment)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Deploy:", oldK8sObj.Name, ", new Deploy:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenIngressEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.ExtensionsV1beta1().RESTClient(), "ingresses", "", fields.Everything())
	extrasMap := map[string]string{"type":"ingress", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&v1beta1.Ingress{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ingress := obj.(*v1beta1.Ingress)
				if _, ok := ingress.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   ingress,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Ingress: ", ingress.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				ingress := obj.(*v1beta1.Ingress)
				if _, ok := ingress.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   ingress,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Ingress: ", ingress.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*v1beta1.Ingress)
				newK8sObj := newObj.(*v1beta1.Ingress)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Ingress:", oldK8sObj.Name, ", new Ingress:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenNetworkPolicyEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.NetworkingV1().RESTClient(), "networkpolicies", "", fields.Everything())
	extrasMap := map[string]string{"type":"networkPolicy", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&networkingV1.NetworkPolicy{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				networkPolicy := obj.(*networkingV1.NetworkPolicy)
				if _, ok := networkPolicy.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   networkPolicy,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add NetworkPolicy: ", networkPolicy.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				networkPolicy := obj.(*networkingV1.NetworkPolicy)
				if _, ok := networkPolicy.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   networkPolicy,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete NetworkPolicy: ", networkPolicy.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*networkingV1.NetworkPolicy)
				newK8sObj := newObj.(*networkingV1.NetworkPolicy)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old NetworkPolicy:", oldK8sObj.Name, ", new NetworkPolicy:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenClusterRoleBindingEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.RbacV1().RESTClient(), "clusterrolebindings", "", fields.Everything())
	extrasMap := map[string]string{"type":"clusterRoleBinding", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&rbacV1.ClusterRoleBinding{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				clusterRoleBinding := obj.(*rbacV1.ClusterRoleBinding)
				if _, ok := clusterRoleBinding.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   clusterRoleBinding,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add ClusterRoleBinding: ", clusterRoleBinding.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				clusterRoleBinding := obj.(*rbacV1.ClusterRoleBinding)
				if _, ok := clusterRoleBinding.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   clusterRoleBinding,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete ClusterRoleBinding: ", clusterRoleBinding.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*rbacV1.ClusterRoleBinding)
				newK8sObj := newObj.(*rbacV1.ClusterRoleBinding)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old ClusterRoleBinding:", oldK8sObj.Name, ", new ClusterRoleBinding:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenClusterRoleEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.RbacV1().RESTClient(), "clusterroles", "", fields.Everything())
	extrasMap := map[string]string{"type":"clusterRole", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&rbacV1.ClusterRole{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				clusterRole := obj.(*rbacV1.ClusterRole)
				if _, ok := clusterRole.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   clusterRole,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add ClusterRole: ", clusterRole.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				clusterRole := obj.(*rbacV1.ClusterRole)
				if _, ok := clusterRole.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   clusterRole,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete ClusterRole: ", clusterRole.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*rbacV1.ClusterRole)
				newK8sObj := newObj.(*rbacV1.ClusterRole)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old ClusterRole:", oldK8sObj.Name, ", new ClusterRole:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenRoleBindingEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.RbacV1().RESTClient(), "rolebindings", "", fields.Everything())
	extrasMap := map[string]string{"type":"roleBinding", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&rbacV1.RoleBinding{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				roleBinding := obj.(*rbacV1.RoleBinding)
				if _, ok := roleBinding.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   roleBinding,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add RoleBinding: ", roleBinding.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				roleBinding := obj.(*rbacV1.RoleBinding)
				if _, ok := roleBinding.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   roleBinding,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete RoleBinding: ", roleBinding.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*rbacV1.RoleBinding)
				newK8sObj := newObj.(*rbacV1.RoleBinding)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old RoleBinding:", oldK8sObj.Name, ", new RoleBinding:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenRoleEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.RbacV1().RESTClient(), "roles", "", fields.Everything())
	extrasMap := map[string]string{"type":"role", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&rbacV1.Role{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				role := obj.(*rbacV1.Role)
				if _, ok := role.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   role,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Role: ", role.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				role := obj.(*rbacV1.Role)
				if _, ok := role.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   role,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Role: ", role.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*rbacV1.Role)
				newK8sObj := newObj.(*rbacV1.Role)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Role:", oldK8sObj.Name, ", new Role:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenServiceAccountEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "serviceAccounts", "", fields.Everything())
	extrasMap := map[string]string{"type":"serviceAccount", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.ServiceAccount{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				serviceAccount := obj.(*coreV1.ServiceAccount)
				if _, ok := serviceAccount.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   serviceAccount,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add ServiceAccount: ", serviceAccount.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				serviceAccount := obj.(*coreV1.ServiceAccount)
				if _, ok := serviceAccount.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   serviceAccount,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete ServiceAccount: ", serviceAccount.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.ServiceAccount)
				newK8sObj := newObj.(*coreV1.ServiceAccount)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old ServiceAccount:", oldK8sObj.Name, ", new ServiceAccount:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenSecretEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "secrets", "", fields.Everything())
	extrasMap := map[string]string{"type":"secret", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.Secret{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				secret := obj.(*coreV1.Secret)
				if _, ok := secret.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   secret,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add Secret: ", secret.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				secret := obj.(*coreV1.Secret)
				if _, ok := secret.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   secret,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete Secret: ", secret.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.Secret)
				newK8sObj := newObj.(*coreV1.Secret)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old Secret:", oldK8sObj.Name, ", new Secret:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenConfigMapEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "configMaps", "", fields.Everything())
	extrasMap := map[string]string{"type":"configMap", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.ConfigMap{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				configMap := obj.(*coreV1.ConfigMap)
				if _, ok := configMap.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   configMap,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add ConfigMap: ", configMap.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				configMap := obj.(*coreV1.ConfigMap)
				if _, ok := configMap.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   configMap,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete ConfigMap: ", configMap.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.ConfigMap)
				newK8sObj := newObj.(*coreV1.ConfigMap)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old ConfigMap:", oldK8sObj.Name, ", new ConfigMap:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenPVCEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "persistentVolumeClaims", "", fields.Everything())
	extrasMap := map[string]string{"type":"pvc", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.PersistentVolumeClaim{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				persistentVolumeClaim := obj.(*coreV1.PersistentVolumeClaim)
				if _, ok := persistentVolumeClaim.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   persistentVolumeClaim,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add PVC: ", persistentVolumeClaim.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				persistentVolumeClaim := obj.(*coreV1.PersistentVolumeClaim)
				if _, ok := persistentVolumeClaim.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   persistentVolumeClaim,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete PVC: ", persistentVolumeClaim.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.PersistentVolumeClaim)
				newK8sObj := newObj.(*coreV1.PersistentVolumeClaim)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old PVC:", oldK8sObj.Name, ", new PVC:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenPVEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.CoreV1().RESTClient(), "persistentVolumes", "", fields.Everything())
	extrasMap := map[string]string{"type":"persistentVolume", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&coreV1.PersistentVolume{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				persistentVolume := obj.(*coreV1.PersistentVolume)
				if _, ok := persistentVolume.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   persistentVolume,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add PV: ", persistentVolume.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				persistentVolume := obj.(*coreV1.PersistentVolume)
				if _, ok := persistentVolume.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   persistentVolume,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete PV: ", persistentVolume.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*coreV1.PersistentVolume)
				newK8sObj := newObj.(*coreV1.PersistentVolume)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old PV:", oldK8sObj.Name, ", new PV:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenDaemonSetEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.AppsV1().RESTClient(), "daemonSets", "", fields.Everything())
	extrasMap := map[string]string{"type":"daemonSet", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&appsV1.DaemonSet{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				daemonSet := obj.(*appsV1.DaemonSet)
				if _, ok := daemonSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   daemonSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add DaemonSet: ", daemonSet.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				daemonSet := obj.(*appsV1.DaemonSet)
				if _, ok := daemonSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   daemonSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete DaemonSet: ", daemonSet.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*appsV1.DaemonSet)
				newK8sObj := newObj.(*appsV1.DaemonSet)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old DaemonSet:", oldK8sObj.Name, ", new DaemonSet:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenReplicaSetEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.AppsV1().RESTClient(), "replicaSets", "", fields.Everything())
	extrasMap := map[string]string{"type":"replicaSet", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&appsV1.ReplicaSet{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				replicaSet := obj.(*appsV1.ReplicaSet)
				if _, ok := replicaSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   replicaSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add ReplicaSet: ", replicaSet.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				replicaSet := obj.(*appsV1.ReplicaSet)
				if _, ok := replicaSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   replicaSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete ReplicaSet: ", replicaSet.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*appsV1.ReplicaSet)
				newK8sObj := newObj.(*appsV1.ReplicaSet)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old ReplicaSet:", oldK8sObj.Name, ", new ReplicaSet:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) ListenStateFullSetSetEvents() (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(k.kcs.AppsV1().RESTClient(), "statefulSets", "", fields.Everything())
	extrasMap := map[string]string{"type":"statefulSet", "agent":config.AgentName}
	return cache.NewInformer(
		watchlist,
		&appsV1.StatefulSet{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				statefulSet := obj.(*appsV1.StatefulSet)
				if _, ok := statefulSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   statefulSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.ADD,
							Extras: extrasMap,
						},
					})
					log.Println("add StatefulSet: ", statefulSet.Name)
				}
			},
			DeleteFunc: func(obj interface{}) {
				statefulSet := obj.(*appsV1.StatefulSet)
				if _, ok := statefulSet.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   statefulSet,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.DELETE,
							Extras: extrasMap,
						},
					})
					log.Println("delete StatefulSet: ", statefulSet.Name)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldK8sObj := oldObj.(*appsV1.StatefulSet)
				newK8sObj := newObj.(*appsV1.StatefulSet)

				type KubeObject struct {
					oldK8sObj, newK8sObj interface{}
				}
				obj := KubeObject{
					oldK8sObj: oldK8sObj,
					newK8sObj: newK8sObj,
				}
				if _, ok := newK8sObj.Labels["klovercloud_ci"]; ok {
					k.kubeEventPublisher.Publish(v1.KubeEventMessage{
						Body:   obj,
						Header: v1.MessageHeader{
							Offset:  0,
							Command: enums.UPDATE,
							Extras: extrasMap,
						},
					})
					log.Println("old StatefulSet:", oldK8sObj.Name, ", new StatefulSet:", newK8sObj.Name)
				}
			},
		},
	)
}

func (k k8sService) Apply(data unstructured.Unstructured) error {
	_, err := k.Deploy(&data)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

func (k k8sService) Deploy(data *unstructured.Unstructured) (bool, error) {
	for {
		version := data.GetAPIVersion()
		kind := data.GetKind()
		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		apiResourceList, err := k.discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return false, err
		}
		apiResources := apiResourceList.APIResources
		var resource *metaV1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}
		if resource == nil {
			return false, fmt.Errorf("unknown resource kind: %s", kind)
		}

		groupVersionResource := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}
		namespace := "default"
		if strings.Compare(data.GetNamespace(), "_all") == 0 {
			namespace = data.GetNamespace()
		}

		if resource.Namespaced {
			_, err = k.dynamicClient.Resource(groupVersionResource).Namespace(namespace).Create(context.Background(), data, metaV1.CreateOptions{})

		} else {
			_, err = k.dynamicClient.Resource(groupVersionResource).Create(context.Background(), data, metaV1.CreateOptions{})
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func (k k8sService) UpdateDeployment(resource v1.Resource) error {
	subject := v1.Subject{resource.Step, "Initiating  deployment ...", resource.Name, resource.Namespace, resource.ProcessId, nil, nil, resource.Pipeline}
	subject.EventData = make(map[string]interface{})
	subject.EventData["footmark"] = enums.UPDATE_RESOURCE
	subject.EventData["log"] = subject.Log
	subject.EventData["reason"] = "n/a"
	subject.EventData["status"] = enums.INITIALIZING
	subject.EventData["claim"] = strconv.Itoa(resource.Claim)
	go k.notifyAll(subject)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetDeployment(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: ", getErr)
			subject.Log = "Failed to get latest version of Deployment: " + getErr.Error()
			subject.EventData["log"] = subject.Log
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.DEPLOYMENT_FAILED
			go k.notifyAll(subject)
			return getErr
		}
		for i, each := range resource.Images {
			if i > len(result.Spec.Template.Spec.Containers)-1 {
				subject.Log = "index out of bound! ignoring container for " + each
				subject.EventData["log"] = subject.Log
				subject.EventData["footmark"] = enums.UPDATE_RESOURCE
				subject.EventData["status"] = enums.PROCESSING
				go k.notifyAll(subject)
			} else {
				result.Spec.Template.Spec.Containers[i].Image = each
			}
		}

		listOptions := metaV1.ListOptions{LabelSelector: labels.FormatLabels(result.Labels)}
		podList, err := k.kcs.CoreV1().Pods(resource.Namespace).List(context.TODO(), listOptions)
		if err != nil {
			subject.Log = "Failed to list Existing Pods!"
			subject.EventData["log"] = err.Error()
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.PROCESSING
			go k.notifyAll(subject)
		}
		existingPodName := make(map[string]bool)

		for _, each := range podList.Items {
			existingPodName[each.Name] = true
		}
		prev, _ := k.GetDeployment(resource.Name, resource.Namespace)
		if result.Labels == nil {
			result.Labels = make(map[string]string)
		}
		result.Labels["company"] = resource.Pipeline.MetaData.CompanyId
		result.Labels["klovercloud_ci"] = "enabled"
		deploy, updateErr := k.PatchDeploymentObject(prev, result)
		if updateErr != nil {
			log.Println("patchError:", err.Error())
		}
		if updateErr == nil && *deploy.Spec.Replicas > 0 {
			subject.Log = "Waiting until pod is ready!"
			subject.EventData["log"] = subject.Log
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.PROCESSING
			go k.notifyAll(subject)
			var timeout = 30
			err := k.WaitForPodBySelectorUntilRunning(resource.Step, deploy.Name, deploy.Namespace, resource.ProcessId, labels.FormatLabels(deploy.Labels), timeout, existingPodName, resource.Claim)
			if err != nil {
				return err
			}
		}
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}

	subject.Log = "Updated Successfully"
	subject.EventData["log"] = subject.Log
	subject.EventData["footmark"] = enums.POST_AGENT_JOB
	subject.EventData["status"] = enums.SUCCESSFUL
	go k.notifyAll(subject)
	return nil
}

func (k k8sService) PatchDeploymentObject(cur, mod *appsV1.Deployment) (*appsV1.Deployment, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}
	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, err
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, appsV1.Deployment{})
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	out, err := k.kcs.AppsV1().Deployments(cur.Namespace).Patch(context.TODO(), cur.Name, types.StrategicMergePatchType, patch, metaV1.PatchOptions{})
	return out, err
}

func (k k8sService) UpdatePod(resource v1.Resource) error {
	subject := v1.Subject{resource.Step, "", resource.Name, resource.Namespace, resource.ProcessId, map[string]interface{}{"footmark": enums.UPDATE_RESOURCE, "log": "Initiating  deployment ...", "reason": "n/a"}, nil, resource.Pipeline}
	subject.EventData["status"] = enums.INITIALIZING
	go k.notifyAll(subject)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetPod(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Pod: ", getErr)
			subject.Log = "Failed to get latest version of Pod: " + getErr.Error()
			subject.EventData["log"] = subject.Log
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.DEPLOYMENT_FAILED
			go k.notifyAll(subject)
			return getErr
		}

		for i, each := range resource.Images {
			result.Spec.Containers[i].Image = each
		}

		_, updateErr := k.kcs.CoreV1().Pods(resource.Namespace).Update(context.TODO(), result, metaV1.UpdateOptions{})
		k.kcs.CoreV1().Pods(resource.Namespace).Watch(context.TODO(), metaV1.ListOptions{})

		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) UpdateStatefulSet(resource v1.Resource) error {
	subject := v1.Subject{resource.Step, "Initiating  deployment ...", resource.Name, resource.Namespace, resource.ProcessId, nil, nil, resource.Pipeline}
	subject.EventData = make(map[string]interface{})
	subject.EventData["footmark"] = enums.UPDATE_RESOURCE
	subject.EventData["log"] = subject.Log
	subject.EventData["reason"] = "n/a"
	subject.EventData["status"] = enums.INITIALIZING
	subject.EventData["claim"] = strconv.Itoa(resource.Claim)
	go k.notifyAll(subject)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetStatefulSet(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of StatefulSet: ", getErr)
			subject.Log = "Failed to get latest version of StatefulSet: " + getErr.Error()
			subject.EventData["log"] = subject.Log
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.DEPLOYMENT_FAILED
			go k.notifyAll(subject)
			return getErr
		}
		for i, each := range resource.Images {
			if i > len(result.Spec.Template.Spec.Containers)-1 {
				subject.Log = "index out of bound! ignoring container for " + each
				subject.EventData["log"] = subject.Log
				subject.EventData["footmark"] = enums.UPDATE_RESOURCE
				subject.EventData["status"] = enums.PROCESSING
				go k.notifyAll(subject)
			} else {
				result.Spec.Template.Spec.Containers[i].Image = each
			}
		}
		if result.Labels == nil {
			result.Labels = make(map[string]string)
		}
		result.Labels["company"] = resource.Pipeline.MetaData.CompanyId
		result.Labels["klovercloud_ci"] = "enabled"
		listOptions := metaV1.ListOptions{LabelSelector: labels.FormatLabels(result.Labels)}
		podList, err := k.kcs.CoreV1().Pods(resource.Namespace).List(context.TODO(), listOptions)
		if err != nil {
			subject.Log = "Failed to list Existing Pods!"
			subject.EventData["log"] = err.Error()
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.PROCESSING
			go k.notifyAll(subject)
		}
		existingPodRevisions := make(map[string]bool)
		for _, each := range podList.Items {
			existingPodRevisions[each.GetResourceVersion()] = true
		}
		statefulSet, updateErr := k.kcs.AppsV1().StatefulSets(resource.Namespace).Update(context.TODO(), result, metaV1.UpdateOptions{})
		if updateErr == nil && *statefulSet.Spec.Replicas > 0 {
			subject.Log = "Waiting until pod is ready!"
			subject.EventData["log"] = subject.Log
			subject.EventData["footmark"] = enums.POST_AGENT_JOB
			subject.EventData["status"] = enums.PROCESSING
			go k.notifyAll(subject)
			var timeout = 30
			err := k.WaitForPodBySelectorAndRevisionUntilRunning(resource.Step, statefulSet.Name, statefulSet.Namespace, resource.ProcessId, labels.FormatLabels(statefulSet.Labels), timeout, existingPodRevisions, resource.Claim)
			if err != nil {
				return err
			}
		}
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	subject.Log = "Updated Successfully"
	subject.EventData["log"] = subject.Log
	subject.EventData["footmark"] = enums.POST_AGENT_JOB
	subject.EventData["status"] = enums.SUCCESSFUL
	go k.notifyAll(subject)
	return nil
}

func (k k8sService) UpdateDaemonSet(resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetDaemonSet(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		for i, each := range resource.Images {
			result.Spec.Template.Spec.Containers[i].Image = each
		}
		if result.Labels == nil {
			result.Labels = make(map[string]string)
		}
		result.Labels["company"] = resource.Pipeline.MetaData.CompanyId
		result.Labels["klovercloud_ci"] = "enabled"
		_, updateErr := k.kcs.AppsV1().DaemonSets(resource.Namespace).Update(context.TODO(), result, metaV1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) GetDeployment(name, namespace string) (*appsV1.Deployment, error) {
	return k.kcs.AppsV1().Deployments(namespace).Get(context.Background(), name, metaV1.GetOptions{})
}

func (k k8sService) GetPod(name, namespace string) (*coreV1.Pod, error) {
	return k.kcs.CoreV1().Pods(namespace).Get(context.Background(), name, metaV1.GetOptions{})
}

func (k k8sService) GetStatefulSet(name, namespace string) (*appsV1.StatefulSet, error) {
	return k.kcs.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metaV1.GetOptions{})
}

func (k k8sService) GetDaemonSet(name, namespace string) (*appsV1.DaemonSet, error) {
	return k.kcs.AppsV1().DaemonSets(namespace).Get(context.Background(), name, metaV1.GetOptions{})
}

func (k k8sService) isPodRunning(podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".") // progress bar!

		pod, err := k.kcs.CoreV1().Pods(namespace).Get(context.Background(), podName, metaV1.GetOptions{})
		if err != nil {
			return false, err
		}

		for _, each := range pod.Status.ContainerStatuses {
			if each.State.Waiting != nil {
				if each.State.Waiting.Reason == "ImagePullBackOff" {
					return true, errors.New("Pod has error: ImagePullBackOff")
				} else if each.State.Waiting.Reason == "CrashLoopBackOff" {
					return true, errors.New("Pod has error: CrashLoopBackOff")
				}
			}
		}

		switch pod.Status.Phase {
		case coreV1.PodRunning:
			return true, nil
		case coreV1.PodFailed, coreV1.PodSucceeded:
			return false, errors.New("Pod has error!")
		}
		return false, nil
	}
}

func (k k8sService) waitForPodRunning(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, k.isPodRunning(podName, namespace))
}

func (k k8sService) ListNewPodsByRevision(retryCount int, namespace, selector string, revision map[string]bool) (*coreV1.PodList, error) {
	listOptions := metaV1.ListOptions{LabelSelector: selector}
	podList, err := k.kcs.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	var newPods []string

	for _, each := range podList.Items {
		if _, ok := revision[each.GetResourceVersion()]; ok {
			continue
		}
		newPods = append(newPods, each.Name)
	}
	if len(newPods) == 0 && retryCount < 10 {
		time.Sleep(time.Second * 2)
		retryCount = retryCount + 1
		return k.ListNewPods(retryCount, namespace, selector, revision)
	}
	return podList, nil
}

func (k k8sService) ListNewPods(retryCount int, namespace, selector string, existingPodMap map[string]bool) (*coreV1.PodList, error) {
	listOptions := metaV1.ListOptions{LabelSelector: selector}
	podList, err := k.kcs.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	var newPods []string

	for _, each := range podList.Items {
		if _, ok := existingPodMap[each.Name]; ok {
			continue
		}
		newPods = append(newPods, each.Name)
	}
	if len(newPods) == 0 && retryCount < 10 {
		time.Sleep(time.Second * 2)
		retryCount = retryCount + 1
		return k.ListNewPods(retryCount, namespace, selector, existingPodMap)
	}
	return podList, nil
}

func (k k8sService) WaitForPodBySelectorAndRevisionUntilRunning(step, name, namespace, processId, selector string, timeout int, revisions map[string]bool, claim int) error {
	subject := v1.Subject{step, "Listing Pods ...", name, namespace, processId, nil, nil, nil}
	subject.EventData = make(map[string]interface{})
	subject.EventData["log"] = subject.Log
	subject.EventData["footmark"] = enums.POST_AGENT_JOB
	subject.EventData["status"] = enums.PROCESSING
	subject.EventData["reason"] = "n/a"
	subject.EventData["claim"] = strconv.Itoa(claim)
	go k.notifyAll(subject)
	podList, err := k.ListNewPodsByRevision(0, namespace, selector, revisions)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return fmt.Errorf("no pods in %s with selector %s", namespace, selector)
	}
	for _, pod := range podList.Items {
		subject.Log = "Waiting until pods are ready ..."
		subject.EventData["log"] = subject.Log
		go k.notifyAll(subject)
		if err := k.waitForPodRunning(namespace, pod.Name, time.Duration(timeout)*time.Second); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (k k8sService) WaitForPodBySelectorUntilRunning(step, name, namespace, processId, selector string, timeout int, existingPods map[string]bool, claim int) error {
	subject := v1.Subject{step, "Listing Pods ...", name, namespace, processId, nil, nil, nil}
	subject.EventData = make(map[string]interface{})
	subject.EventData["log"] = subject.Log
	subject.EventData["footmark"] = enums.POST_AGENT_JOB
	subject.EventData["status"] = enums.PROCESSING
	subject.EventData["reason"] = "n/a"
	subject.EventData["claim"] = strconv.Itoa(claim)
	go k.notifyAll(subject)
	podList, err := k.ListNewPods(0, namespace, selector, existingPods)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return fmt.Errorf("no pods in %s with selector %s", namespace, selector)
	}

	for _, pod := range podList.Items {
		subject.Log = "Waiting until pods are ready ..."
		subject.EventData["log"] = subject.Log
		go k.notifyAll(subject)
		if err := k.waitForPodRunning(namespace, pod.Name, time.Duration(timeout)*time.Second); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (k k8sService) notifyAll(subject v1.Subject) {
	for _, observer := range k.observerList {
		go observer.Listen(subject)
	}
}

// NewK8sService returns K8s type service.
func NewK8sService(Kcs *kubernetes.Clientset, dynamicClient dynamic.Interface, discoveryClient *discovery.DiscoveryClient, observerList []service.Observer, kubeEventPublisher service.KubeEventPublisher) service.K8s {
	return &k8sService{
		kcs:                Kcs,
		dynamicClient:      dynamicClient,
		discoveryClient:    discoveryClient,
		observerList:       observerList,
		kubeEventPublisher: kubeEventPublisher,
	}
}
