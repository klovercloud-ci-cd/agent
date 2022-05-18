package service

type KubeEvent interface {
	GetK8sObjectChangeEvents()
}
