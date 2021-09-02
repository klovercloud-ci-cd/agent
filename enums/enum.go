package enums
type RESOURCE_TYPE string
const  (
	DEPLOYMENT= RESOURCE_TYPE("deployment")
	STATEFULSET= RESOURCE_TYPE("statefulset")
	DAEMONSET= RESOURCE_TYPE("daemonset")
	POD= RESOURCE_TYPE("pod")
	REPLICASET= RESOURCE_TYPE("replicaset")
)
