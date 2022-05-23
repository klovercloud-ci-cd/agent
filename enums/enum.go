package enums

// RESOURCE_TYPE pipeline resource types
type RESOURCE_TYPE string

const ( // CERTIFICATE k8s certificate as resource
	CERTIFICATE = RESOURCE_TYPE("certificate")
	// CLUSTER_ROLE k8s cluster_role as resource
	CLUSTER_ROLE = RESOURCE_TYPE("clusterRole")
	// CLUSTER_ROLE_BINDGING k8s cluster_role_binding as resource
	CLUSTER_ROLE_BINDGING = RESOURCE_TYPE("clusterRoleBinding")
	// CONFIG_MAP k8s config_map as resource
	CONFIG_MAP = RESOURCE_TYPE("configMap")
	// DAEMONSET k8s daemonset as resource
	DAEMONSET = RESOURCE_TYPE("daemonset")
	// DEPLOYMENT k8s deployment as resource
	DEPLOYMENT = RESOURCE_TYPE("deployment")
	// INGRESS k8s ingress as resource
	INGRESS = RESOURCE_TYPE("ingress")
	// NAMESPACE k8s namespace as resource
	NAMESPACE = RESOURCE_TYPE("namespace")
	// NETWORK_POLICY k8s network_policy as resource
	NETWORK_POLICY = RESOURCE_TYPE("networkPolicy")
	// NODE k8s node as resource
	NODE = RESOURCE_TYPE("node")
	// PERSISTENT_VOLUME k8s persistent_volume as resource
	PERSISTENT_VOLUME = RESOURCE_TYPE("persistentVolume")
	// PERSISTENT_VOLUME_CLAIM k8s persistent_volume_claim as resource
	PERSISTENT_VOLUME_CLAIM = RESOURCE_TYPE("persistentVolumeClaim")
	// POD k8s pod as resource
	POD = RESOURCE_TYPE("pod")
	// REPLICASET k8s replicaset as resource
	REPLICASET = RESOURCE_TYPE("replicaset")
	// ROLE k8s role as resource
	ROLE = RESOURCE_TYPE("role")
	// ROLE_BINDING k8s role_binding as resource
	ROLE_BINDING = RESOURCE_TYPE("roleBinding")
	// SECRET k8s secret as resource
	SECRET = RESOURCE_TYPE("secret")
	// SERVICE k8s service as resource
	SERVICE = RESOURCE_TYPE("service")
	// SERVICE_ACCOUNT k8s service_account as resource
	SERVICE_ACCOUNT = RESOURCE_TYPE("serviceAccount")
	// STATEFULSET k8s statefulset as resource
	STATEFULSET = RESOURCE_TYPE("statefulset")
	// EVENT k8s statefulset as resource
	EVENT = RESOURCE_TYPE("event")
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

// ENVIRONMENT run environment
type ENVIRONMENT string

const (
	// PRODUCTION production environment
	PRODUCTION = ENVIRONMENT("PRODUCTION")
	// DEVELOP development environment
	DEVELOP = ENVIRONMENT("DEVELOP")
	// TEST test environment
	TEST = ENVIRONMENT("TEST")
)

// FOOTMARK process footmark (step breakdown)
type FOOTMARK string

const (
	// INIT_AGNET_JOB FOOTMARK name
	INIT_AGNET_JOB = FOOTMARK("init_agent_job")
	// POST_AGENT_JOB FOOTMARK name
	POST_AGENT_JOB = FOOTMARK("post_agent_job")
	// UPDATE_RESOURCE FOOTMARK name
	UPDATE_RESOURCE = FOOTMARK("update_resource")
)

// Command kafka command
type Command string

const (
	// Kube object ADD command
	ADD = Command("ADD")
	// Kube object UPDATE command
	UPDATE = Command("UPDATE")
	// Kube object DELETE command
	DELETE = Command("DELETE")
)
