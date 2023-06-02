package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/pkg/term"
	"github.com/spf13/cobra"
)

var ttyDevice string

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
	if ttyDevice != "" {
		return term.Open(ttyDevice, term.Speed(9600), term.RawMode)
	}

	for _, devNode := range []string{
		"/dev/ttyFAN0",
		"/dev/ttyUSB0",
	} {
		// If we found a usable tty previously, skip the probing.
		fd, err := term.Open(devNode, term.Speed(9600), term.RawMode)
		if err == nil {
			ttyDevice = devNode
			return fd, nil
		}

		switch {
		case errors.Is(err, fs.ErrNotExist):
		case errors.Is(err, fs.ErrPermission):
			fmt.Fprintf(os.Stderr, "The user does not have rights to the serial port device (`/dev/ttyUSB0`)\n")
			fmt.Fprintf(os.Stderr, "- Screw around with udev rules\n")
			fmt.Fprintf(os.Stderr, "- Add your user to the appropriate group (`uucp`, `dialout`)\n")
			os.Exit(1)
		default:
			return nil, fmt.Errorf("failed to open serial port: %w", err)
		}
	}

	fmt.Fprintf(os.Stderr, "The serial port device (`/dev/(ttyFAN0,ttyUSB0)`) does not appear to exist.\n")
	fmt.Fprintf(os.Stderr, " 1. Add `dtoverlay=dwc2,dr_mode=host` to `/boot/config.txt`\n")
	fmt.Fprintf(os.Stderr, " 2. Fully reboot the device (unplug the power)\n")
	os.Exit(1)

	panic("BUG: unreached")
}
