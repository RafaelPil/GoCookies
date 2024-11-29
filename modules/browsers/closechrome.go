package browsers

import (
	"os/exec"
	"log"
)

func CloseChrome() {
	cmd := exec.Command("taskkill", "/IM", "chrome.exe", "/F")
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to close Chrome: %v", err)
	} else {
		log.Println("Chrome closed successfully.")
	}
}
