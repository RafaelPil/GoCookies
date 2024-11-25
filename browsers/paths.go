package browsers

import (
	"os"
	"path/filepath"
)

// GetChromiumBrowsers returns a map of Chromium-based browsers and their user data directories.
func GetChromiumBrowsers() map[string]string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		// Handle error if the home directory cannot be determined
		panic("Unable to determine the user's home directory")
	}

	// Dynamically construct full paths based on the user's home directory
	return map[string]string{
		"Chromium":             filepath.Join(userHome, "AppData", "Local", "Chromium", "User Data"),
		"Thorium":              filepath.Join(userHome, "AppData", "Local", "Thorium", "User Data"),
		"Chrome":               filepath.Join(userHome, "AppData", "Local", "Google", "Chrome", "User Data", "Default"),
		"Chrome (x86)":         filepath.Join(userHome, "AppData", "Local", "Google(x86)", "Chrome", "User Data", "Default"),
		"Chrome SxS":           filepath.Join(userHome, "AppData", "Local", "Google", "Chrome SxS", "User Data", "Default"),
		"Maple":                filepath.Join(userHome, "AppData", "Local", "MapleStudio", "ChromePlus", "User Data"),
		"Iridium":              filepath.Join(userHome, "AppData", "Local", "Iridium", "User Data"),
		"7Star":                filepath.Join(userHome, "AppData", "Local", "7Star", "7Star", "User Data"),
		"CentBrowser":          filepath.Join(userHome, "AppData", "Local", "CentBrowser", "User Data"),
		"Chedot":               filepath.Join(userHome, "AppData", "Local", "Chedot", "User Data"),
		"Vivaldi":              filepath.Join(userHome, "AppData", "Local", "Vivaldi", "User Data"),
		"Kometa":               filepath.Join(userHome, "AppData", "Local", "Kometa", "User Data"),
		"Elements":             filepath.Join(userHome, "AppData", "Local", "Elements Browser", "User Data"),
		"Epic Privacy Browser": filepath.Join(userHome, "AppData", "Local", "Epic Privacy Browser", "User Data"),
		"Uran":                 filepath.Join(userHome, "AppData", "Local", "uCozMedia", "Uran", "User Data"),
		"Fenrir":               filepath.Join(userHome, "AppData", "Local", "Fenrir Inc", "Sleipnir5", "setting", "modules", "ChromiumViewer"),
		"Catalina":             filepath.Join(userHome, "AppData", "Local", "CatalinaGroup", "Citrio", "User Data"),
		"Coowon":               filepath.Join(userHome, "AppData", "Local", "Coowon", "Coowon", "User Data"),
		"Liebao":               filepath.Join(userHome, "AppData", "Local", "liebao", "User Data"),
		"QIP Surf":             filepath.Join(userHome, "AppData", "Local", "QIP Surf", "User Data"),
		"Orbitum":              filepath.Join(userHome, "AppData", "Local", "Orbitum", "User Data"),
		"Dragon":               filepath.Join(userHome, "AppData", "Local", "Comodo", "Dragon", "User Data"),
		"360Browser":           filepath.Join(userHome, "AppData", "Local", "360Browser", "Browser", "User Data"),
		"Maxthon":              filepath.Join(userHome, "AppData", "Local", "Maxthon3", "User Data"),
		"K-Melon":              filepath.Join(userHome, "AppData", "Local", "K-Melon", "User Data"),
		"CocCoc":               filepath.Join(userHome, "AppData", "Local", "CocCoc", "Browser", "User Data"),
		"Brave":                filepath.Join(userHome, "AppData", "Local", "BraveSoftware", "Brave-Browser", "User Data"),
		"Amigo":                filepath.Join(userHome, "AppData", "Local", "Amigo", "User Data"),
		"Torch":                filepath.Join(userHome, "AppData", "Local", "Torch", "User Data"),
		"Sputnik":              filepath.Join(userHome, "AppData", "Local", "Sputnik", "Sputnik", "User Data"),
		"Edge":                 filepath.Join(userHome, "AppData", "Local", "Microsoft", "Edge", "User Data"),
		"DCBrowser":            filepath.Join(userHome, "AppData", "Local", "DCBrowser", "User Data"),
		"Yandex":               filepath.Join(userHome, "AppData", "Local", "Yandex", "YandexBrowser", "User Data"),
		"UR Browser":           filepath.Join(userHome, "AppData", "Local", "UR Browser", "User Data"),
		"Slimjet":              filepath.Join(userHome, "AppData", "Local", "Slimjet", "User Data"),
		"Opera":                filepath.Join(userHome, "AppData", "Roaming", "Opera Software", "Opera Stable"),
		"OperaGX":              filepath.Join(userHome, "AppData", "Roaming", "Opera Software", "Opera GX Stable"),
	}
}
