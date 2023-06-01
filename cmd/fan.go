package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	cmdFanOff = "pwm_000"
	cmdFanOn  = "pwm_100"
)

var speedCmds = map[int]string{
	0:   cmdFanOff,
	25:  "pwm_025",
	50:  "pwm_050",
	75:  "pwm_075",
	100: cmdFanOn,
}

var fanCmd = &cobra.Command{
	Use:       "fan {0 25 50 75 100}",
	Short:     "Set the fan speed",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"0", "25", "50", "75", "100"},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetSpeed, _ := strconv.Atoi(args[0])
		cmdStr := speedCmds[targetSpeed]

		fd, err := getCtrlTty()
		if err != nil {
			return fmt.Errorf("deskpi-ctl/fan: %w", err)
		}
		defer fd.Close()

		if _, err = fd.Write([]byte(cmdStr)); err != nil {
			return fmt.Errorf("deskpi-ctl/fan: failed to send command: %w", err)
		}

		return nil
	},
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
}

var fanDaemonCmd = &cobra.Command{
	Use:   "fan-daemon",
	Short: "Turn the fan on/off based on CPU temp",
	RunE: func(cmd *cobra.Command, args []string) error {
		fd, err := getCtrlTty()
		if err != nil {
			return fmt.Errorf("deskpi-ctl/fan-daemon: %w", err)
		}
		defer fd.Close()

		var fanIsOn bool
		for {
			temp, err := getCPUTemp()
			if err != nil {
				return err
			}

			var toWrite []byte

			switch {
			case temp < 45:
				if fanIsOn {
					toWrite = []byte(cmdFanOff)
					fanIsOn = false
				}
			case temp > 50:
				if !fanIsOn {
					toWrite = []byte(cmdFanOn)
					fanIsOn = true
				}
			}

			if toWrite != nil {
				if _, err = fd.Write(toWrite); err != nil {
					return fmt.Errorf("deskpi-ctl/fan-daemon: failed to send command: %w", err)
				}
			}
			time.Sleep(time.Second)
		}
	},
}

func getCPUTemp() (int, error) {
	const tempSensorFn = "/sys/class/thermal/thermal_zone0/temp"

	b, err := os.ReadFile(tempSensorFn)
	if err != nil {
		return 0, fmt.Errorf("deskpi-ctl/fan-daemon: failed to poll CPU temp: %w", err)
	}

	temp, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return 0, fmt.Errorf("deskpi-ctl/fan-daemon: failed to parse CPU temp: %w", err)
	}

	return temp / 1000, nil
}

func init() {
	rootCmd.AddCommand(fanCmd)
	rootCmd.AddCommand(fanDaemonCmd)
}
