package internal

const cpu_usage = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=\"%s\"}) by (pod)"

//get the number of running pods group by namespace
const active_pods_by_namespace = "count(kube_pod_info) by (namespace)"

//get all runnnig pods
const all_active_pods = "count(kube_pod_info)"

//Workload Memory working
const memory_allocation = "sum(container_memory_working_set_bytes{namespace=\"%s\"}) by (pod)"

var queryDict = map[string]string{
	"cpu_usage":                cpu_usage,
	"all_active_pods":          all_active_pods,
	"active_pods_by_namespace": active_pods_by_namespace,
	"memory_allocation":        memory_allocation,
}
