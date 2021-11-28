package internal

const cpu_usage = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=\"%s\"}) by (pod)"

//get the number of running pods group by namespace
const active_pods_by_namespace = "count(kube_pod_info) by (namespace)"

//get all runnnig pods
const all_active_pods = "count(kube_pod_info)"

//Workload Memory working
const memory_allocation = "sum(container_memory_working_set_bytes{namespace=\"%s\"}) by (pod)"

//Memory usage
const memory_usage = "container_memory_working_set_bytes{pod_name=~\"compute-.*\", image!=\"\", container_name!=\"POD\"}"

//CPU allocation
const cpu_allocation = "avg(kube_pod_container_resource_limits_cpu_cores{pod=~\"compute-.*\"})"

/******* CO2 EMISSION FORMULA block ******/

const cloud_provider_usage = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate)"

//calculate cloud energy conversion factors [kWh]
const min_watts_coeficient = 0.71
const max_watts_coeficient = 4.26
const avg_cpu_utilization = 0.5 //50%

const avg_watts = min_watts_coeficient + avg_cpu_utilization*(max_watts_coeficient-min_watts_coeficient)

//cpu utilization per hour
const cpu_utilization_per_hour = "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate) / 3600"
const compute_watts_hours = avg_watts * avg_cpu_utilization

//power usage effectiveness coeficient
const pue_coeficient = 1.1

//grid emission factors coeficient for europe-west1 zone
const emission_factors_coeficient = 0.000196

//co2 formula but it has problem with string and float
//const co2_emission = cloud_provider_usage * compute_watts_hours * pue_coeficient * emission_factors_coeficient

//co2 formula
const co2_emission = "(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate)) * (0.71 + 0.5 * (4.26-0.71)) * (sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate) / 3600) * 1.1 * 0.000196"

var queryDict = map[string]string{
	"cpu_usage":                cpu_usage,
	"all_active_pods":          all_active_pods,
	"active_pods_by_namespace": active_pods_by_namespace,
	"memory_allocation":        memory_allocation,
	"memory_usage":             memory_usage,
	"cpu_allocation":           cpu_allocation,
	"co2_emission":             co2_emission,
}
