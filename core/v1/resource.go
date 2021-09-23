package v1

import "github.com/klovercloud-ci/enums"

type Resource struct {
	Step string `json:"step"`
	ProcessId string `json:"process_id"`
	Descriptors *[]interface{} `json:"descriptors" yaml:"descriptors"`
	Type     enums.RESOURCE_TYPE `json:"type"`
	Name string                  `json:"name"`
	Namespace string             `json:"namespace"`
	Replica int32                `json:"replica"`
	Images [] struct {
		ImageIndex int `json:"image_index"`
		Image      string `json:"image"`
	}`json:"images"`
}