package Service

import (
	"os"
	"strconv"
	"strings"
)

const loadDeviceArgument = "--load-device="

// GetStartupDeviceMode returns the driver mode requested by an elevated
// restart. -1 means this is a normal application start.
func (g *AppMain) GetStartupDeviceMode() int {
	for _, argument := range os.Args[1:] {
		if !strings.HasPrefix(argument, loadDeviceArgument) {
			continue
		}
		mode, err := strconv.Atoi(strings.TrimPrefix(argument, loadDeviceArgument))
		if err == nil && mode >= 0 && mode <= 2 {
			return mode
		}
	}
	return -1
}
