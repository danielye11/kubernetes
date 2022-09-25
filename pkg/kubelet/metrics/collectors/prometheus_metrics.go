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
	"k8s.io/klog/v2"
	statsapi "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
)

var (
	promCollectorDesc = metrics.NewDesc(
		"danielyecollector",
		"Bytes used by the container's logs on the filesystem.",
		[]string{
			"danielyecollector",
			"danielyecollector",
			"danielyecollector",
			"danielyecollector",
		}, nil,
		metrics.ALPHA,
		"",
	)
)

type PromCollector struct {
	metrics.BaseStableCollector

	podStats func() ([]statsapi.PodStats, error)
}

// Check if logMetricsCollector implements necessary interface
var _ metrics.StableCollector = &logMetricsCollector{}

// NewLogMetricsCollector implements the metrics.StableCollector interface and
// exposes metrics about container's log volume size.
func NewPromCollector(podStats func() ([]statsapi.PodStats, error)) metrics.StableCollector {
	klog.Infof("danielyeprometheuslogger is being hit")
	return &PromCollector{
		podStats: podStats,
	}
}

// DescribeWithStability implements the metrics.StableCollector interface.
func (c *PromCollector) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- promCollectorDesc
}

// CollectWithStability implements the metrics.StableCollector interface.
func (c *PromCollector) CollectWithStability(ch chan<- metrics.Metric) {
	podStats, err := c.podStats()
	if err != nil {
		klog.ErrorS(err, "Failed to get pod stats")
		return
	}

	for _, ps := range podStats {
		for _, c := range ps.Containers {
			if c.Logs != nil && c.Logs.UsedBytes != nil {
				ch <- metrics.NewLazyConstMetric(
					promCollectorDesc,
					metrics.GaugeValue,
					float64(*c.Logs.UsedBytes),
					ps.PodRef.UID,
					ps.PodRef.Namespace,
					ps.PodRef.Name,
					c.Name,
				)
			}
		}
	}
}
