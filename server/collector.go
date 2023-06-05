package server

import (
	"errors"
	"strconv"
	"strings"
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
	instructionInfo *prometheus.Desc
	cycleInfo       *prometheus.Desc
}

func NewPerfEventCollector() (*PerfEventCollector, error) {
	return &PerfEventCollector{
		instructionInfo: prometheus.NewDesc(
			prometheus.BuildFQName("perf_event", "instruction", "count"),
			"Perf event instruction count",
			[]string{"id"}, nil,
		),
		cycleInfo: prometheus.NewDesc(
			prometheus.BuildFQName("perf_event", "cycle", "count"),
			"Perf event cycle count",
			[]string{"id"}, nil,
		),
	}, nil
}

func (c *PerfEventCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.instructionInfo
	ch <- c.cycleInfo
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
			c.instructionInfo,
			prometheus.CounterValue,
			float64(totalInstructions),
			id,
		)

		ch <- prometheus.MustNewConstMetric(
			c.cycleInfo,
			prometheus.CounterValue,
			float64(totalcycles),
			id,
		)

	}

}

func AddCgroupPerfEventCollector(config *PerfEventConfig) error {
	var cpu []int

	if config.Cpuset != "" {
		str := strings.Split(config.Cpuset, "-")
		if len(str) < 2 {
			return errors.New("wrong cpu set format")
		}
		lo, err := strconv.Atoi(str[0])
		if err != nil {
			return errors.New("wrong cpu set format")
		}

		hi, err := strconv.Atoi(str[1])
		if err != nil {
			return errors.New("wrong cpu set format")
		}

		if hi < lo || lo < 0 || hi > utils.NumCPU {
			return errors.New("wrong cpu set")
		}

		cpu = make([]int, hi-lo+1)

		for i := lo; i <= hi; i++ {
			cpu[i-lo] = i
		}
	}

	hwProfilers, err := utils.NewCgroupPerfeventProfilerMap(config.Cgroup, cpu)
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
