package v1

import (
	"github.com/klovercloud-ci-cd/agent/enums"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Resource agent applicable workload info.
type Resource struct {
	Step        string                       `json:"step"`
	ProcessId   string                       `json:"process_id"`
	Descriptors *[]unstructured.Unstructured `json:"descriptors" yaml:"descriptors"`
	Type        enums.RESOURCE_TYPE          `json:"type"`
	Name        string                       `json:"name"`
	Namespace   string                       `json:"namespace"`
	Replica     int32                        `json:"replica"`
	Images      []string                     `json:"images"`
	Pipeline    *Pipeline                    `bson:"pipeline" json:"pipeline"`
	Claim       int                          `bson:"claim" json:"claim"`
}

// Pipeline pipeline stuct
type Pipeline struct {
	MetaData   PipelineMetadata `json:"_metadata" yaml:"_metadata"`
	ApiVersion string           `json:"api_version" yaml:"api_version"`
	Name       string           `json:"name"  yaml:"name"`
	ProcessId  string           `json:"process_id" yaml:"process_id"`
	Steps      []Step           `json:"steps" yaml:"steps"`
}

// Step pipeline step.
type Step struct {
	Name string   `json:"name" yaml:"name"`
	Next []string `json:"next" yaml:"next"`
}

// PipelineMetadata pipeline metadata
type PipelineMetadata struct {
	CompanyId string `json:"company_id" yaml:"company_id"`
}
