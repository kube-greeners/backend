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
const cpu_allocation = "sum(namespace_cpu:kube_pod_container_resource_requests:sum)"

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

//co2 formula
var co2_emission = fmt.Sprintf("(%s) * (%f * (%s)) * %f * %f", cpu_usage, avg_watts, cpu_utilization_per_hour, pue_coeficient, emission_factors_coeficient)

var queryDict = map[string]string{
	"cpu_usage":         cpu_usage,
	"all_active_pods":   all_active_pods,
	"memory_allocation": memory_allocation,
	"memory_usage":      memory_usage,
	"cpu_allocation":    cpu_allocation,
	"co2_emission":      co2_emission,
}
