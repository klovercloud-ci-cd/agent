package v1

// Subject subject that observers listen
type Subject struct {
	Step, Log                  string
	Name, Namespace, ProcessId string
	EventData                  map[string]interface{}
	ProcessLabel               map[string]string
	Pipeline                   *Pipeline
}
