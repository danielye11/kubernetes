/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collectors

import (
	"k8s.io/component-base/metrics"
	runtimeapi "k8s.io/cri-api/pkg/apis"
	kubeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
)

var (
	descContainerMetrics = metrics.NewDesc(
		"kubelet_container_metrics_kubernetes_test",
		"Container metrics prometheus test.",
		[]string{
			"label_key",
			"label_value",
		}, nil,
		metrics.ALPHA,
		"",
	)
)

type containerMetricsCollector struct {
	metrics.BaseStableCollector

	// listContainerStatsFunc func(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error)
	// runtimeapi.ContainerStatsManager.ListContainerStats
	// *runtimeapi.ContainerStatsManagerfunc(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error)
	manager runtimeapi.ContainerStatsManager

	// stats runtimeapi.RuntimeService.containerStatsManager
}

// Check if logMetricsCollector implements necessary interface
var _ metrics.StableCollector = &containerMetricsCollector{}

// NewLogMetricsCollector implements the metrics.StableCollector interface and
// exposes metrics about container's log volume size.
func NewContainerMetricsCollector(manager runtimeapi.ContainerStatsManager) metrics.StableCollector {
	return &containerMetricsCollector{
		manager: manager,
	}
}

// DescribeWithStability implements the metrics.StableCollector interface.
func (c *containerMetricsCollector) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- descContainerMetrics
}

// CollectWithStability implements the metrics.StableCollector interface.
func (c *containerMetricsCollector) CollectWithStability(ch chan<- metrics.Metric) {

	// var m []kubeapi.Metric
	// m[0].Label = &kubeapi.LabelPair{
	// 	Name:  "danielye: prometheus dummy metric name",
	// 	Value: "danielye: prometheus dummy metric value",
	// }
	// m[0].Gauge = &kubeapi.Gauge{
	// 	Value: 1.0,
	// }
	// m[0].Type = kubeapi.MetricType_COUNTER
	// m[0].TimestampMs = 17
	// var prometheus_label runtime.LabelPair
	// var label_array []runtime.LabelPair = &runtime.Lba
	// {&runtime.LabelPair{
	// 	Name:  []string{"danielye: prometheus dummy metric name",",1"},
	// 	Value: []string{"danielye: prometheus dummy metric value"},
	// }}
	// m.Label = label_array[]runtime.LabelPair{}
	// var label_name runtime.LabelPair

	cs, err := c.manager.ListContainerStats(&kubeapi.ContainerStatsFilter{})
	if err != nil {
		klog.ErrorS(err, "Failed to get container stats")
		return
	}
	for _, c := range cs {
		if c.PrometheusMetric != nil {
			ch <- metrics.NewLazyConstMetric(
				descContainerMetrics,
				metrics.GaugeValue,
				0,
				c.PrometheusMetric.Label.Name,
				c.PrometheusMetric.Label.Value,
			)
		}
	}
	// resp, errors = &kubeapi.ListContainerStatsResponse{cs, []&kubeapi.Metric{}}

	// podStats, err := c.podStats()
	// if err != nil {
	// 	klog.ErrorS(err, "Failed to get pod stats")
	// 	return
	// }

	// for _, ps := range podStats {
	// 	for _, c := range ps.Containers {
	// 		if c.Logs != nil && c.Logs.UsedBytes != nil {
	// 			ch <- metrics.NewLazyConstMetric(
	// 				descLogSize,
	// 				metrics.GaugeValue,
	// 				float64(*c.Logs.UsedBytes),
	// 				ps.PodRef.UID,
	// 				ps.PodRef.Namespace,
	// 				ps.PodRef.Name,
	// 				c.Name,
	// 			)
	// 		}
	// 	}
	// }
}
