//go:build windows

package Service

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/windows"
)

// RequestDeviceElevation keeps normal startup unelevated and requests UAC only
// when a driver is about to be loaded. On success, the elevated instance is
// given the selected mode and the current instance exits after replying to UI.
func (g *AppMain) RequestDeviceElevation(mode int) string {
	if mode < 0 || mode > 2 {
		return "error:invalid driver mode"
	}
	if windows.GetCurrentProcessToken().IsElevated() {
		return "ready"
	}

	executable, err := os.Executable()
	if err != nil {
		return "error:" + err.Error()
	}
	executable, err = filepath.Abs(executable)
	if err != nil {
		return "error:" + err.Error()
	}

	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(executable)
	arguments, _ := windows.UTF16PtrFromString(loadDeviceArgument + strconv.Itoa(mode))
	workingDirectory, _ := windows.UTF16PtrFromString(filepath.Dir(executable))
	if err = windows.ShellExecute(0, verb, file, arguments, workingDirectory, 1); err != nil {
		if errors.Is(err, windows.ERROR_CANCELLED) {
			return "cancelled"
		}
		return "error:" + err.Error()
	}

	go func() {
		time.Sleep(700 * time.Millisecond)
		if app != nil {
			app.Quit()
		}
	}()
	return "restarting"
}
