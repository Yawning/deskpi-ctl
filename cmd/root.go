package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/pkg/term"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "deskpi-ctl - Control a DeskPi Pro",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	DisableFlagsInUseLine: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getCtrlTty() (*term.Term, error) {
	fd, err := term.Open("/dev/ttyUSB0", term.Speed(9600), term.RawMode)
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			fmt.Fprintf(os.Stderr, "The serial port device (`/dev/ttyUSB0`) does not appear to exist.\n")
			fmt.Fprintf(os.Stderr, " 1. Add `dtoverlay=dwc2,dr_mode=host` to `/boot/config.txt`\n")
			fmt.Fprintf(os.Stderr, " 2. Fully reboot the device (unplug the power)\n")
			os.Exit(1)
		case errors.Is(err, fs.ErrPermission):
			fmt.Fprintf(os.Stderr, "The user does not have rights to the serial port device (`/dev/ttyUSB0`)\n")
			fmt.Fprintf(os.Stderr, "- Screw around with udev rules\n")
			fmt.Fprintf(os.Stderr, "- Add your user to the appropriate group (`uucp`, `dialout`)\n")
			os.Exit(1)
		}
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}
	return fd, nil
}
