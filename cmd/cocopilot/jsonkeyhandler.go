package main

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/tmc/keyring"
)

// JSONKeyHandler provides read/write access to a JSON value, such as an JSON
// web token, stored within a persistent keyring provided by the OS.
type JSONKeyHandler[T comparable] struct {
	Service  string
	Username string
}

// Set stores the key value to the underlying keyring.
func (handler JSONKeyHandler[T]) Set(value *T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal key value: %w", err)
	}

	err = keyring.Set(handler.Service, handler.Username, string(data))
	if err != nil {
		return fmt.Errorf("failed to store key: %w", err)
	}

	return nil
}

// Get retrieves the key value from the underlying keyring.
func (handler JSONKeyHandler[T]) Get() (*T, error) {
	data, err := keyring.Get(handler.Service, handler.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key: %w", err)
	}

	value := new(T)

	err = json.NewDecoder(strings.NewReader(data)).Decode(&value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key value: %w", err)
	}

	return value, nil
}
