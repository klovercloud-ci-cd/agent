package v1

type Subject struct {
	Step,Log string
	Name,Namespace,ProcessId string
	EventData map[string]interface{}
	ProcessLabel map[string]string
}
