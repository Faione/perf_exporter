package server

import (
	"sync"

	"github.com/Faione/perf_exporter/utils"
	"github.com/hodgesds/perf-utils"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	collectorsLock sync.RWMutex
	collectors     = make(map[string]map[int]*perf.HardwareProfiler)
)

type PerfEventCollector struct {
	cgroupPerfEventInfo *prometheus.Desc
}

func NewPerfEventCollector() (*PerfEventCollector, error) {
	return &PerfEventCollector{
		cgroupPerfEventInfo: prometheus.NewDesc(
			prometheus.BuildFQName("cgroup", "perf_event", "count"),
			"Perf event count on cgroup",
			[]string{"event", "id"}, nil,
		),
	}, nil
}

func (c *PerfEventCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cgroupPerfEventInfo
}

func (c *PerfEventCollector) Collect(ch chan<- prometheus.Metric) {
	collectorsLock.RLock()
	defer collectorsLock.RUnlock()

	for id, profilers := range collectors {

		var totalInstructions, totalcycles uint64

		for _, hwProf := range profilers {

			hwProfile := &perf.HardwareProfile{}
			if err := (*hwProf).Profile(hwProfile); err != nil {
				break
			}
			totalInstructions += *hwProfile.Instructions
			totalcycles += *hwProfile.CPUCycles

		}

		ch <- prometheus.MustNewConstMetric(
			c.cgroupPerfEventInfo,
			prometheus.CounterValue,
			float64(totalInstructions),
			"Instructions", id,
		)

		ch <- prometheus.MustNewConstMetric(
			c.cgroupPerfEventInfo,
			prometheus.CounterValue,
			float64(totalcycles),
			"cycles", id,
		)

	}

}

func AddCgroupPerfEventCollector(config *PerfEventConfig) error {
	hwProfilers, err := utils.NewCgroupPerfeventProfilerMap(config.Cgroup, nil)
	if err != nil {
		return err
	}

	collectorsLock.Lock()
	defer collectorsLock.Unlock()
	collectors[config.Cgroup] = hwProfilers
	return nil
}

func DelCgroupPerfEventCollector(config *PerfEventConfig) error {
	collectorsLock.Lock()
	defer collectorsLock.Unlock()
	delete(collectors, config.Cgroup)
	return nil
}
