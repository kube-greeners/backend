package internal

import (
	"fmt"
)

//get all running pods
const all_active_pods = "count(kube_pod_info)"

// Ozan: sum((container_memory_usage_bytes{namespace!=""} - container_memory_working_set_bytes{namespace!=""}) / 1024 / 1024) by (namespace)
// If you look, usage_bytes are bigger than working_set_bytes

//Workload Memory working
const memory_allocation = "sum(container_memory_usage_bytes{namespace!=\"\"})"

//Memory usage
const memory_usage = "sum(container_memory_working_set_bytes{namespace!=\"\"})"

//CPU allocation
const cpu_allocation = "avg(kube_pod_container_resource_limits_cpu_cores{})"

const cpu_usage = "sum(rate(container_cpu_usage_seconds_total[6h]))"

//calculate cloud energy conversion factors [kWh]
// https://github.com/cloud-carbon-footprint/cloud-carbon-footprint/blob/trunk/packages/gcp/src/domain/GcpFootprintEstimationConstants.ts
const min_watts_coeficient = 0.71
const max_watts_coeficient = 4.26
const avg_cpu_utilization = 0.5 //50%

const avg_watts = min_watts_coeficient + avg_cpu_utilization*(max_watts_coeficient-min_watts_coeficient)

// Utilization of CPU time for the last week summed up
const cpu_utilization_per_hour = "sum(avg_over_time(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate[7d])) / 3600"

const pue_coeficient = 1.1 // GKE

const emission_factors_coeficient = 0.000196 * 1000 // europe-west1
// Multiply by 1000 because, the factor is in metric tons

//co2 formula but it has problem with string and float
//const co2_emission = cloud_provider_usage * compute_watts_hours * pue_coeficient * emission_factors_coeficient

//co2 formula
var co2_emission = fmt.Sprintf("(%s) * (%f * (%s)) * %f * %f", cpu_usage, avg_watts, cpu_utilization_per_hour, pue_coeficient, emission_factors_coeficient)

// const co2_emission = "(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate)) * (0.71 + 0.5 * (4.26-0.71)) * (sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate) / 3600) * 1.1 * 0.000196"

func decreaseAtCertainHour(startHour int, endHour int, query string, scaleCoef float32) string {
	newQuery := fmt.Sprintf(
		"(%s and (absent(hour() < %d) or absent(hour() > %d))) or (%s * %f and (absent(hour() > %d) or absent(hour() < %d)))",
		query, startHour, endHour, query, scaleCoef, startHour, endHour)
	return newQuery
}

var queryDict = map[string]string{
	"cpu_usage":         cpu_usage,
	"all_active_pods":   all_active_pods,
	"memory_allocation": memory_allocation,
	"memory_usage":      memory_usage,
	"cpu_allocation":    cpu_allocation,
	"co2_emission":      co2_emission,
	"saved_cpu_at_night": decreaseAtCertainHour(17, 6, queryCPU, 0.7),
	"saved_memory_at_night": decreaseAtCertainHour(17, 6, queryMemory, 0.8),
}
