package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Power off the unit",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		fd, err := getCtrlTty()
		if err != nil {
			return fmt.Errorf("deskpi-ctl/shutdown: %w", err)
		}
		defer fd.Close()

		if _, err = fd.Write([]byte("power_off")); err != nil {
			return fmt.Errorf("deskpi-ctl/shutdown: failed to send command: %w", err)
		}

		return nil
	},
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
}

func init() {
	rootCmd.AddCommand(shutdownCmd)
}
