package logic

import (
	"context"
	"fmt"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/enums"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"log"
	"strings"
)

type k8sService struct {
	kcs             *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	discoveryClient *discovery.DiscoveryClient
	observerList    []service.Observer
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
	subject := v1.Subject{resource.Step, "", resource.Name, resource.Namespace, resource.ProcessId, map[string]interface{}{"log": "Initiating  deployment ...", "reason": "n/a"}, nil}
	subject.EventData["status"] = enums.INITIALIZING
	go k.notifyAll(subject)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetDeployment(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: ", getErr)
			subject.Log = "Failed to get latest version of Deployment: " + getErr.Error()
			subject.EventData["log"] = subject.Log
			subject.EventData["status"] = enums.DEPLOYMENT_FAILED
			go k.notifyAll(subject)
			return getErr
		}
		result.Spec.Replicas = &resource.Replica
		for i, each := range resource.Images {
			if i > len(result.Spec.Template.Spec.Containers)-1 {
				subject.Log = "index out of bound! ignoring container for " + each
				subject.EventData["log"] = subject.Log
				subject.EventData["status"] = enums.PROCESSING
				go k.notifyAll(subject)
			} else {
				result.Spec.Template.Spec.Containers[i].Image = each
			}
		}
		_, updateErr := k.kcs.AppsV1().Deployments(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		subject.Log = "Update failed: " + retryErr.Error()
		subject.EventData["log"] = subject.Log
		subject.EventData["status"] = enums.DEPLOYMENT_FAILED
		go k.notifyAll(subject)
		return retryErr
	}

	subject.Log = "Updated Successfully"
	subject.EventData["log"] = subject.Log
	subject.EventData["status"] = enums.SUCCESSFUL
	go k.notifyAll(subject)
	return nil
}

func (k k8sService) UpdatePod(resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetPod(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		for i, each := range resource.Images {
			result.Spec.Containers[i].Image = each
		}
		_, updateErr := k.kcs.CoreV1().Pods(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) UpdateStatefulSet(resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetStatefulSet(resource.Name, resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		result.Spec.Replicas = &resource.Replica
		for i, each := range resource.Images {
			result.Spec.Template.Spec.Containers[i].Image = each
		}
		_, updateErr := k.kcs.AppsV1().StatefulSets(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
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
		_, updateErr := k.kcs.AppsV1().DaemonSets(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) GetDeployment(name, namespace string) (*apiV1.Deployment, error) {
	return k.kcs.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k k8sService) GetPod(name, namespace string) (*coreV1.Pod, error) {
	return k.kcs.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k k8sService) GetStatefulSet(name, namespace string) (*apiV1.StatefulSet, error) {
	return k.kcs.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k k8sService) GetDaemonSet(name, namespace string) (*apiV1.DaemonSet, error) {
	return k.kcs.AppsV1().DaemonSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k k8sService) notifyAll(subject v1.Subject) {
	for _, observer := range k.observerList {
		go observer.Listen(subject)
	}
}

// NewK8sService returns K8s type service.
func NewK8sService(Kcs *kubernetes.Clientset, dynamicClient dynamic.Interface, discoveryClient *discovery.DiscoveryClient, observerList []service.Observer) service.K8s {
	return &k8sService{
		kcs:             Kcs,
		dynamicClient:   dynamicClient,
		discoveryClient: discoveryClient,
		observerList:    observerList,
	}
}
