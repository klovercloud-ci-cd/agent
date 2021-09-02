package logic

import (
	"context"
	"fmt"
	"github.com/klovercloud-ci/core/v1/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	v1 "github.com/klovercloud-ci/core/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

type k8sService struct {
	Kcs   *kubernetes.Clientset
	observerList []service.Observer
}

func (k k8sService) UpdateDeployment(deployment apiV1.Deployment, resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.Kcs.AppsV1().Deployments(deployment.Namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Replicas = &resource.Replica
		for _,each:=range resource.Images{
			result.Spec.Template.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) UpdatePod(pod coreV1.Pod, resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.Kcs.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
		}
		for _,each:=range resource.Images{
			result.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.CoreV1().Pods(pod.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) UpdateStatefulSet(statefulSet apiV1.StatefulSet, resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.Kcs.AppsV1().StatefulSets(statefulSet.Namespace).Get(context.TODO(), statefulSet.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of StatefulSet: %v", getErr))
		}
		result.Spec.Replicas = &resource.Replica
		for _,each:=range resource.Images{
			result.Spec.Template.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.AppsV1().StatefulSets(statefulSet.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Println("Update failed: %v", retryErr)
		return retryErr
	}
	return nil
}

func (k k8sService) UpdateDaemonSet(daemonSet apiV1.DaemonSet, resource v1.Resource) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := k.Kcs.AppsV1().DaemonSets(daemonSet.Namespace).Get(context.TODO(), daemonSet.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of DaemonSet: %v", getErr))
		}
		for _,each:=range resource.Images{
			result.Spec.Template.Spec.Containers[each.ImageIndex].Image=each.Image
		}
		_, updateErr := k.Kcs.AppsV1().DaemonSets(daemonSet.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
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