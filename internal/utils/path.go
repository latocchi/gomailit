package utils

import (
	"log"
	"os"
	"path/filepath"
)

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.Mode().IsRegular()
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}

func getAppConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Unable to find user config directory: %v", err)
	}

	appDir := filepath.Join(configDir, "gomailit")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		log.Fatalf("Unable to create config directory: %v", err)
	}

	return appDir
}

func TokenPath() string {
	return filepath.Join(getAppConfigDir(), "token.json")
}
