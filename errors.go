package main

import (
	"errors"
)

var (
	ErrNotFound = errors.New("Requested file not found")
)
