package logic

import (
	"context"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"log"
)

type k8sService struct {
	Kcs   *kubernetes.Clientset
	observerList []service.Observer
}


func (k k8sService) UpdateDeployment( resource v1.Resource) error {
	subject:=v1.Subject{resource.Step,"",resource.Name,resource.Namespace,resource.ProcessId,map[string]interface{}{"log":"Initiating  deployment ...","reason":"n/a"},nil}
	subject.EventData["status"]=enums.INITIALIZING
	go k.notifyAll(subject)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetDeployment(resource.Name,resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: ", getErr)
			subject.Log="Failed to get latest version of Deployment: "+ getErr.Error()
			subject.EventData["log"]=subject.Log
			subject.EventData["status"]=enums.FAILED
			go k.notifyAll(subject)
	    	return getErr
		}
		result.Spec.Replicas = &resource.Replica
		for i,each:=range resource.Images{
			if i>len(result.Spec.Template.Spec.Containers)-1{
				subject.Log="index out of bound! ignoring container for "+each.Image
				subject.EventData["log"]=subject.Log
				subject.EventData["status"]=enums.PROCESSING
				go k.notifyAll(subject)
			}else {
				result.Spec.Template.Spec.Containers[each.ImageIndex].Image = each.Image
			}
		}
		_, updateErr := k.Kcs.AppsV1().Deployments(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		subject.Log="Update failed: "+ retryErr.Error()
		subject.EventData["log"]=subject.Log
		subject.EventData["status"]=enums.FAILED
		go k.notifyAll(subject)
		return retryErr
	}

	subject.Log="Updated Successfully"
	subject.EventData["log"]=subject.Log
	subject.EventData["status"]=enums.SUCCESSFUL
	go k.notifyAll(subject)
	return nil
}

func (k k8sService) UpdatePod(resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.GetPod(resource.Name,resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		for _,each:=range resource.Images{
			result.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.CoreV1().Pods(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
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
		result, getErr := k.GetStatefulSet(resource.Name,resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		result.Spec.Replicas = &resource.Replica
		for _,each:=range resource.Images{
			result.Spec.Template.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.AppsV1().StatefulSets(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
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
		result, getErr := k.GetDaemonSet(resource.Name,resource.Namespace)
		if getErr != nil {
			log.Println("Failed to get latest version of Deployment: %v", getErr)
			return getErr
		}
		for _,each:=range resource.Images{
			result.Spec.Template.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.AppsV1().DaemonSets(resource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) GetDeployment(name, namespace string) (*apiV1.Deployment, error) {
	return k.Kcs.AppsV1().Deployments(namespace).Get(context.Background(),name,metav1.GetOptions{})
}

func (k k8sService) GetPod(name, namespace string) (*coreV1.Pod, error) {
	return k.Kcs.CoreV1().Pods(namespace).Get(context.Background(),name,metav1.GetOptions{})
}

func (k k8sService) GetStatefulSet(name, namespace string) (*apiV1.StatefulSet, error) {
	return k.Kcs.AppsV1().StatefulSets(namespace).Get(context.Background(),name,metav1.GetOptions{})
}

func (k k8sService) GetDaemonSet(name,namespace string) (*apiV1.DaemonSet, error) {
	return k.Kcs.AppsV1().DaemonSets(namespace).Get(context.Background(),name,metav1.GetOptions{})
}

func (k8s k8sService)notifyAll(subject v1.Subject){
	for _, observer := range k8s.observerList {
		go observer.Listen(subject)
	}
}

func NewK8sService(Kcs *kubernetes.Clientset,observerList []service.Observer) service.K8s {
	return &k8sService{
		Kcs:                 Kcs,
		observerList: observerList,
	}
}