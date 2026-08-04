//go:build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// Do not resolve system DLLs from the executable/current directory. Files
	// such as dwmapi.dll are commonly left beside downloaded applications and
	// would otherwise shadow the genuine copy in System32 during Wails startup.
	const flags = windows.LOAD_LIBRARY_SEARCH_SYSTEM32 | windows.LOAD_LIBRARY_SEARCH_USER_DIRS
	if err := windows.SetDefaultDllDirectories(flags); err != nil {
		fmt.Fprintln(os.Stderr, "configure DLL search path:", err)
	}
}
