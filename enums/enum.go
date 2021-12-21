package enums

// RESOURCE_TYPE pipeline resource types
type RESOURCE_TYPE string

const (
	// DEPLOYMENT k8s deployment as resource
	DEPLOYMENT = RESOURCE_TYPE("deployment")
	// STATEFULSET k8s statefulset as resource
	STATEFULSET = RESOURCE_TYPE("statefulset")
	// DAEMONSET k8s daemonset as resource
	DAEMONSET = RESOURCE_TYPE("daemonset")
	// POD k8s pod as resource
	POD = RESOURCE_TYPE("pod")
	// REPLICASET k8s replicaset as resource
	REPLICASET = RESOURCE_TYPE("replicaset")
)

// PIPELINE_STATUS pipeline status
type PIPELINE_STATUS string

const (
	// DEPLOYMENT_FAILED step deploy has been FAILED
	DEPLOYMENT_FAILED = PIPELINE_STATUS("FAILED")
	// PROCESSING step deploy has been PROCESSING
	PROCESSING = PIPELINE_STATUS("PROCESSING")
	// TERMINATING step deploy has been TERMINATING
	TERMINATING = PIPELINE_STATUS("TERMINATING")
	// INITIALIZING step deploy has been INITIALIZING
	INITIALIZING = PIPELINE_STATUS("INITIALIZING")
	// SUCCESSFUL step deploy has been SUCCESSFUL
	SUCCESSFUL = PIPELINE_STATUS("SUCCESSFUL")
	// ERROR step deploy has been ERROR
	ERROR = PIPELINE_STATUS("ERROR")
)

// PROCESS_STATUS pipeline steps status
type PROCESS_STATUS string

const (
	// COMPLETED pipeline steps status completed
	COMPLETED = PROCESS_STATUS("completed")
	// FAILED pipeline steps status failed
	FAILED = PROCESS_STATUS("failed")
	// COMPLETED pipeline steps status completed
	PAUSED = PROCESS_STATUS("paused")
)
