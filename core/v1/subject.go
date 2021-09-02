package v1

type Subject struct {
	Name,Namespace string
	EventData map[string]interface{}
	ProcessLabel map[string]string
}
