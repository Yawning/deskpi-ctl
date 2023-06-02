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
	cmdFanOff    = cmdFanPrefix + "000"
	cmdFanOn     = cmdFanPrefix + "100"
	cmdFanPrefix = "pwm_"
)

var fanCmd = &cobra.Command{
	Use:   "fan {0 ... 100}",
	Short: "Set the fan speed",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("deskpi-ctl/fan: expected integer fan speed")
		}
		if i < 0 || i > 100 {
			return fmt.Errorf("deskpi-ctl/fan: expected fan speed in the range [0, 100]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetSpeed, _ := strconv.Atoi(args[0])

		cmdStr := fmt.Sprintf("%s%03d", cmdFanPrefix, targetSpeed)

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
