package main

import (
	"fmt"
	"strings"
)

func validate(path string) error {
	if len(path) < 2 {
		return fmt.Errorf("Path too short")
	}
	if strings.Contains(path[1:], "/") {
		return fmt.Errorf("Path may not contain '/'")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("Path may not contain '..'")
	}
	return nil
}
