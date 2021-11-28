package internal

import "fmt"

const queryCPU = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=\"%s\"})"

//get running pods group by namespace
const queryPodsByNamespace = "count(kube_pod_info) by (namespace)"

//get all runnnig pods
const queryAllPods = "count(kube_pod_info)"

//Workload Memory Utilization
const queryMemory = "sum(container_memory_working_set_bytes{namespace=\"%s\"})"

func decreaseAtCertainHour(startHour int, endHour int, query string, scaleCoef float32) string {
	newQuery := fmt.Sprintf(
		"(%s and (absent(hour() < %d) or absent(hour() > %d))) or (%s * %f and (absent(hour() > %d) or absent(hour() < %d)))",
		query, startHour, endHour, query, scaleCoef, startHour, endHour)
	return newQuery
}

var queryDict = map[string]string{
	"cpu":                queryCPU,
	"pods":               queryAllPods,
	"podsByNamespace":    queryPodsByNamespace,
	"memory":             queryMemory,
	"saved_cpu_at_night": decreaseAtCertainHour(17, 6, queryCPU, 0.7),
	"saved_memory_at_night": decreaseAtCertainHour(17, 6, queryMemory, 0.8),
}
