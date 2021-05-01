package main

import (
	"fmt"
	"strings"
)

func validate(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("Invalid request. Path contains '..'")
	}
	return nil
}
