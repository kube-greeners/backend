package internal

const queryCPU = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=\"%s\"}) by (pod)"

//get running pods group by namespace
const queryPodsByNamespace = "count(kube_pod_info) by (namespace)"

//get all runnnig pods
const queryAllPods = "count(kube_pod_info)"

//Workload Memory Utilization
const queryMemory = "sum(container_memory_working_set_bytes{namespace=\"%s\"}) by (pod)"

var queryDict = map[string]string{
	"cpu":             queryCPU,
	"pods":            queryAllPods,
	"podsByNamespace": queryPodsByNamespace,
	"memory":          queryMemory,
}
