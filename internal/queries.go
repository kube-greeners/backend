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

const cpu_usage = "sum(rate(container_cpu_usage_seconds_total{namespace=~\"%s\"}[6h]))"

// return 1 when kube-green is not running, empty otherwise
const kg_not_running = "absent(max(kube_green_replicas_sleeping)>0)"

// number of hours kube-green was not running over the past week (out of 168)
const number_hours_kg_not_running_over_1w = "sum_over_time(absent(max(kube_green_replicas_sleeping)>0)[1w:1h])"

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

const memory_coefficient = 0.000392 // GKE

//co2 formula
var cpu_co2_emission = fmt.Sprintf("(%s) * (%f * (%s)) * %f * %f", cpu_usage, avg_watts, cpu_utilization_per_hour, pue_coeficient, emission_factors_coeficient)

var memory_co2_emission = fmt.Sprintf("(%s) * %f", memory_usage, memory_coefficient)

var co2_emission = fmt.Sprintf("(%s + %s) * %f * %f", cpu_co2_emission, memory_co2_emission, pue_coeficient, emission_factors_coeficient)

///////////////// Co2 calculation when kube-green is not running

func getResourceAmountWithoutKG(resource string) string {
	return fmt.Sprintf("(%s) * (%s)", resource, kg_not_running)
}

var cpu_co2_emission_no_kg = fmt.Sprintf("(%s) * (%f * (%s)) * %f * %f", getResourceAmountWithoutKG(cpu_usage), avg_watts, getResourceAmountWithoutKG(cpu_utilization_per_hour), pue_coeficient, emission_factors_coeficient)

var memory_co2_emission_no_kg = fmt.Sprintf("(%s) * %f", getResourceAmountWithoutKG(memory_usage), memory_coefficient)

var co2_emission_no_kg = fmt.Sprintf("(%s + %s) * %f * %f", cpu_co2_emission_no_kg, memory_co2_emission_no_kg, pue_coeficient, emission_factors_coeficient)

///////////////////

func sum_over_time_and_step(query string, time string, step string) string {
	return fmt.Sprintf("sum_over_time((%s)[%s:%s])", query, time, step)
}

const namespace_names = "sum(kube_namespace_labels) by (namespace)"

var estimmated_co2_emission_no_kg = fmt.Sprintf("(%s) * 168 / (%s)", sum_over_time_and_step(co2_emission, "1w", "1h"), number_hours_kg_not_running_over_1w)
var saved_co2_emission = fmt.Sprintf("(%s) - (%s)", co2_emission, estimmated_co2_emission_no_kg)

var queryDict = map[string]string{
	"cpu_usage":                    cpu_usage,
	"all_active_pods":              all_active_pods,
	"memory_allocation":            memory_allocation,
	"memory_usage":                 memory_usage,
	"cpu_allocation":               cpu_allocation,
	"co2_emission":                 co2_emission,
	"co2_emission_with_kube_green": co2_emission,
	"saved_co2_emission":           saved_co2_emission,
	"namespace_names":              namespace_names,
}
