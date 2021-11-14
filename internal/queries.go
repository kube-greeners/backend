package internal

const queryCPU = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=\"%s\"}) by (pod)"

var queryDict = map[string]string{
	"cpu": queryCPU,
}
