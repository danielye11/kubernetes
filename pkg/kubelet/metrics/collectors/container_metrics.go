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
	"time"

	"k8s.io/component-base/metrics"
	runtimeapi "k8s.io/cri-api/pkg/apis"
	kubeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
)

var (
	containerUsageCoreNanoSecondsDesc = metrics.NewDesc("container_usage_core_nano_seconds",
		"Cumulative CPU usage (sum across all cores) since object creation",
		[]string{"container"},
		nil,
		metrics.ALPHA,
		"")
	// Total CPU usage (sum of all cores) averaged over the sample window.
	// The "core" unit can be interpreted as CPU core-nanoseconds per second.
	containerUsageNanoCoresDesc = metrics.NewDesc("container_usage_nano_cores",
		"Total CPU usage (sum of all cores) averaged over the sample window, the core unit can be interpreted as CPU core-nanoseconds per second",
		[]string{"container"},
		nil,
		metrics.ALPHA,
		"")
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
	klog.Infof("danielyelogger is being hit")
	return &containerMetricsCollector{
		manager: manager,
	}
}

// DescribeWithStability implements the metrics.StableCollector interface.
func (c *containerMetricsCollector) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- containerUsageCoreNanoSecondsDesc
}

// CollectWithStability implements the metrics.StableCollector interface.
func (mc *containerMetricsCollector) CollectWithStability(ch chan<- metrics.Metric) {

	cs, err := mc.manager.ListContainerStats(&kubeapi.ContainerStatsFilter{})
	if err != nil {
		klog.ErrorS(err, "Failed to get container stats")
		return
	}
	for _, c := range cs {
		mc.collectContainerCPUMetrics(ch, c)
	}

}

func (mc *containerMetricsCollector) collectContainerCPUMetrics(ch chan<- metrics.Metric, s *kubeapi.ContainerStats) {
	if s.Cpu == nil {
		return
	}
	mc.collectUsageCoreNanoSeconds(ch, s)
	mc.collectUsageNanoCores(ch, s)
}

func (c *containerMetricsCollector) collectUsageCoreNanoSeconds(ch chan<- metrics.Metric, s *kubeapi.ContainerStats) {
	if s.Cpu == nil || s.Cpu.UsageCoreNanoSeconds == nil {
		return
	}

	ch <- metrics.NewLazyConstMetric(containerUsageCoreNanoSecondsDesc, metrics.CounterValue, float64(s.Cpu.UsageCoreNanoSeconds.Value)/float64(time.Second), s.Attributes.Id)
}

func (c *containerMetricsCollector) collectUsageNanoCores(ch chan<- metrics.Metric, s *kubeapi.ContainerStats) {
	if s.Cpu == nil || s.Cpu.UsageNanoCores == nil {
		return
	}

	ch <- metrics.NewLazyConstMetric(containerUsageNanoCoresDesc, metrics.CounterValue, float64(s.Cpu.UsageNanoCores.Value), s.Attributes.Id)
}
