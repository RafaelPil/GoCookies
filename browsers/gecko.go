package browsers

import (
	"log"
	"path/filepath"
)

type Gecko struct{}

// GetCookies extracts cookies for Gecko-based browsers
func (g *Gecko) GetCookies(path string) ([]Cookie, error) {
	dbPath := filepath.Join(path, "cookies.sqlite")
	db, err := GetDBConnection(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT name, value, host, path, expiry FROM moz_cookies`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookies []Cookie
	for rows.Next() {
		var name, host, cookiePath string
		var value []byte
		var expiry int64

		if err := rows.Scan(&name, &value, &host, &cookiePath, &expiry); err != nil {
			log.Printf("Failed to scan row: %v", err) // Log the error and continue
			continue
		}

		// Skip empty or invalid cookies
		if name == "" || host == "" || cookiePath == "" || value == nil {
			continue
		}

		// Convert the expiry timestamp to a human-readable format (optional)
		// You could replace this if you want to handle expiry differently
		cookie := Cookie{
			Name:       name,
			Value:      string(value), // No decryption needed
			Host:       host,
			Path:       cookiePath,
			ExpireDate: expiry,
		}

		cookies = append(cookies, cookie)
	}

	return cookies, nil
}
