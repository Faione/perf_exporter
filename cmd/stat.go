package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Faione/perf_exporter/utils"
	"github.com/hodgesds/perf-utils"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	keyInterval = "interval"
	keyCgroup   = "cgroup"
	keyCpuSet   = "cpu_set"
)

func init() {

	statCmd := &cobra.Command{
		Use:   "stat [optons]",
		Short: "stat perf event from pid | cgroup",
		RunE:  stat,
	}

	flags := statCmd.Flags()

	flags.IntP(keyInterval, "i", 0, "Interval on collecting perf event, ms")
	flags.StringP(keyCgroup, "g", "", "Perf from specific cgruop")
	flags.String(keyCpuSet, "", "Cpus on which to collect, if no cpu set, pe will collect on all cpu, cpu can be split by ',' ")

	rootCmd.AddCommand(statCmd)
}

func stat(c *cobra.Command, args []string) error {
	interval, err := c.Flags().GetInt(keyInterval)
	if err != nil {
		return err
	}
	if interval == 0 {
		return fmt.Errorf("interval must be at lease 1")
	}

	cgroup, err := c.Flags().GetString(keyCgroup)
	if err != nil {
		return err
	}

	hwProfilers, err := utils.NewCgroupPerfeventProfilerMap(cgroup, nil)
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGINT, unix.SIGTERM)
	ticker := time.Tick(time.Duration(interval) * time.Millisecond)

	fmt.Printf("%16s%16s\n", "instruciton", "cycle")
LOOP:
	for {
		select {
		case <-sigs:
			break LOOP
		case <-ticker:

			var totalInstructions, totalcycles uint64

			for cpu, hwProf := range hwProfilers {

				hwProfile := &perf.HardwareProfile{}
				if err := (*hwProf).Profile(hwProfile); err != nil {
					return fmt.Errorf("hwProfile on %d failed: %s", cpu, err)

				}
				totalInstructions += *hwProfile.Instructions
				totalcycles += *hwProfile.CPUCycles

			}

			if totalInstructions == 0 {
				fmt.Printf("%16s", "not counted")
			} else {
				fmt.Printf("%16d", totalInstructions)
			}

			if totalcycles == 0 {
				fmt.Printf("%16s\n", "not counted")
			} else {
				fmt.Printf("%16d\n", totalcycles)
			}

		}
	}

	return nil
}
