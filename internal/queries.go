package internal

import (
	"fmt"
)

//get all running pods
const all_active_pods = "count(kube_pod_info{namespace=~\"%s\"})"

//Workload Memory working
const memory_allocation = "sum(container_memory_usage_bytes{namespace=~\"%s\"}) / 1024 / 1024 / 1024" // 1GB = 1024 * 1024 * 1024 Bytes

//Memory usage
const memory_usage = "sum(container_memory_working_set_bytes{namespace=~\"%s\"}) / 1024 / 1024 / 1024"

//CPU allocation
const cpu_allocation = "sum(namespace_cpu:kube_pod_container_resource_requests:sum{namespace=~\"%s\"})"

const cpu_usage = "sum(rate(container_cpu_usage_seconds_total[6h]{namespace=~\"%s\"}))"

//calculate cloud energy conversion factors [kWh]
// https://github.com/cloud-carbon-footprint/cloud-carbon-footprint/blob/trunk/packages/gcp/src/domain/GcpFootprintEstimationConstants.ts
const min_watts_coeficient = 0.71
const max_watts_coeficient = 4.26
const avg_cpu_utilization = 0.5 //50%

const avg_watts = min_watts_coeficient + avg_cpu_utilization*(max_watts_coeficient-min_watts_coeficient)

// Utilization of CPU time for the last week summed up
const cpu_utilization_per_hour = "sum(avg_over_time(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~\"%s\"}[7d])) / 3600"

const pue_coeficient = 1.1 // GKE

const emission_factors_coeficient = 0.000196 * 1000 * 1000 // europe-west1
// Multiply by 1000 * 1000 because the factor is in metric tons and we want grams

//co2 formula
var cpu_co2_emission = fmt.Sprintf("(%s) * (%f * (%s)) * %f * %f", cpu_usage, avg_watts, cpu_utilization_per_hour, pue_coeficient, emission_factors_coeficient)

const memory_coefficient = 0.000392 // GKE
var memory_co2_emission = fmt.Sprintf("(%s) * %f", memory_usage, memory_coefficient)

var co2_emission = fmt.Sprintf("(%s + %s) * %f * %f", cpu_co2_emission, memory_co2_emission, pue_coeficient, emission_factors_coeficient)

func decreaseAtCertainHour(startHour int, endHour int, query string, scaleCoef float32) string {
	newQuery := fmt.Sprintf(
		"(%s * %f and  (absent(hour() < %d) or absent(hour() > %d))) or (%s and (absent(hour() > %d) or absent(hour() < %d)))",
		query, scaleCoef, startHour, endHour, query, startHour, endHour)
	return newQuery
}

var co2_emission_with_kube_green = decreaseAtCertainHour(17, 6, co2_emission, 0.7)
var saved_co2_emissions = fmt.Sprintf("(%s) - (%s)", co2_emission, co2_emission_with_kube_green)

func sum_over_time_and_step(query string, time string, step string) string {
	return fmt.Sprintf("sum_over_time((%s)[%s:%s])", query, time, step)
}

<<<<<<< HEAD
const namespace_names = "sum(kube_namespace_labels) by namespace"
=======
const namespace_names = "sum(kube_namespace_labels) by (namespace)"
>>>>>>> dev

var queryDict = map[string]string{
	"cpu_usage":                    cpu_usage,
	"all_active_pods":              all_active_pods,
	"memory_allocation":            memory_allocation,
	"memory_usage":                 memory_usage,
	"cpu_allocation":               cpu_allocation,
	"co2_emission":                 co2_emission,
	"co2_emission_with_kube_green": co2_emission_with_kube_green,
	"saved_co2_emission":           sum_over_time_and_step(saved_co2_emissions, "1w", "1h"),
	"namespace_names":              namespace_names,
}
