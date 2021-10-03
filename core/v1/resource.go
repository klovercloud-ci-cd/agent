package v1

import (
	"github.com/klovercloud-ci/enums"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resource struct {
	Step string `json:"step"`
	ProcessId string `json:"process_id"`
	Descriptors *[]unstructured.Unstructured  `json:"descriptors" yaml:"descriptors"`
	Type     enums.RESOURCE_TYPE `json:"type"`
	Name string                  `json:"name"`
	Namespace string             `json:"namespace"`
	Replica int32                `json:"replica"`
	Images [] string`json:"images"`
}