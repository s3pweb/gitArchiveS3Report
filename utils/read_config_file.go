package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadConfigFile reads a configuration file and returns a map with key-value pairs
func ReadConfigFile(filepath string) (map[string]string, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Map to store the key-value pairs
	config := make(map[string]string)

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// Split the line on the first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		// Trim spaces around the key and value
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
