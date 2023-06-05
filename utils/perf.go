package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hodgesds/perf-utils"
)

var (
	PERF_FLAG_PID_CGROUP = 1 << 2

	NumCPU = runtime.NumCPU()
)

func NewCgroupPerfeventProfilerMap(cgroup string, cpu []int) (map[int]*perf.HardwareProfiler, error) {
	file, err := os.OpenFile(cgroup, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return nil, fmt.Errorf("open cgroup failed: %s", err)
	}
	fd := file.Fd()

	if len(cpu) == 0 {
		cpu = make([]int, NumCPU)
		for i := 0; i < NumCPU; i++ {
			cpu[i] = i
		}
	}

	hwProfilers := make(map[int]*perf.HardwareProfiler)
	for _, cpu := range cpu {
		hwProf, err := perf.NewHardwareProfiler(int(fd), cpu, perf.CpuInstrProfiler|perf.CpuCyclesProfiler, PERF_FLAG_PID_CGROUP)
		if err != nil {
			return nil, fmt.Errorf("create hwProf on cpu %d failed: %s", cpu, err)
		}

		if err := hwProf.Start(); err != nil {
			return nil, fmt.Errorf("start hwProf on cpu %d failed: %s", cpu, err)

		}

		hwProfilers[cpu] = &hwProf
	}

	return hwProfilers, file.Close()
}
