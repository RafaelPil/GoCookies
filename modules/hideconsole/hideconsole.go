package hideconsole

import (
	"fmt"
	"syscall"
)

var (
    kernel32          = syscall.NewLazyDLL("kernel32.dll")
    user32            = syscall.NewLazyDLL("user32.dll")
    procGetConsoleWnd = kernel32.NewProc("GetConsoleWindow")
    procShowWindow    = user32.NewProc("ShowWindow")
)

const SW_HIDE = 0 // SW_HIDE to hide the window

// Run hides the console window
func Run() {
    // Get console window handle
    hwnd, _, _ := procGetConsoleWnd.Call()
    if hwnd == 0 {
        fmt.Println("Failed to get console window handle")
        return
    }

    // Hide the console window
    _, _, err := procShowWindow.Call(hwnd, SW_HIDE)
    if err != nil {
        fmt.Println("Failed to hide console window:", err)
    } else {
        fmt.Println("Console window hidden successfully.")
    }
}
