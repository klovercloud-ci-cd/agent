package v1

// PipelineProcessEvent Pipeline ProcessEvent struct
type PipelineProcessEvent struct {
	ProcessId string                 `bson:"process_id" json:"process_id"`
	CompanyId string  `bson:"company_id" json:"company_id"`
	Data      map[string]interface{} `bson:"data" json:"data"`
}
