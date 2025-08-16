package command

import (
	"fmt"
	"os"
	"strings"
)

func readHostsFile(hostsFile string) ([]string, error) {
	if hostsFile == "" {
		return []string{}, nil
	}
	
	// Basic path validation
	if strings.Contains(hostsFile, "..") || strings.HasPrefix(hostsFile, "/") {
		return nil, fmt.Errorf("invalid hosts file path: %s", hostsFile)
	}
	
	bytes, err := os.ReadFile(hostsFile)
	if err != nil {
		return nil, fmt.Errorf("could not read input hosts file: %s: %s", hostsFile, err.Error())
	}
	return strings.Split(string(bytes), "\n"), nil
}
