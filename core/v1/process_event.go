package v1
type PipelineProcessEvent struct {
	ProcessId string  `bson:"process_id"`
	Data map[string]interface{}  `bson:"data"`
}
