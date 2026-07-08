// Package state manages persistent and session-scoped state for the application.
//
// This file provides the Global variables store struct, constructor, and all
// CRUD operations (Set, Get, Unset, All, Count) plus command integration.
package state

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hamedhdd/arsenal-ng/internal/config"
)

// =============================================================================
// Global Variables Store
// =============================================================================

// Global holds persistent variables that are saved to disk.
// Variables are set via "set name=value" and can be used in commands.
// Variables persist across different shell sessions via variables.json file.
// Thread-safe for concurrent access.
type Global struct {
	mu        sync.RWMutex
	variables map[string]string
	filePath  string // Path to variables.json
}

// NewGlobal creates a new variables store and loads existing variables from disk.
// If the variables.json file doesn't exist, it will be created on first save.
// Returns an error if the config directory cannot be created or accessed.
func NewGlobal() (*Global, error) {
	filePath, err := config.GetVariablesPath()
	if err != nil {
		log.Printf("ERROR: Failed to get variables path: %v", err)
		return nil, fmt.Errorf("failed to get variables path: %w", err)
	}

	g := &Global{
		variables: make(map[string]string),
		filePath:  filePath,
	}

	// Load existing variables from disk (ignore error if file doesn't exist)
	if err := g.LoadFromFile(); err != nil && !os.IsNotExist(err) {
		log.Printf("ERROR: Failed to load variables from %s: %v", filePath, err)
		return nil, fmt.Errorf("failed to load variables: %w", err)
	}

	if len(g.variables) > 0 {
		log.Printf("Loaded %d variable(s) from %s", len(g.variables), filePath)
	} else {
		log.Printf("No variables found, starting with empty store (file: %s)", filePath)
	}

	return g, nil
}

// =============================================================================
// CRUD Operations
// =============================================================================

// Set stores or updates a variable and saves to disk.
func (g *Global) Set(name, value string) error {
	g.mu.Lock()
	existed := false
	if _, exists := g.variables[name]; exists {
		existed = true
	}
	g.variables[name] = value
	g.mu.Unlock()

	action := "Set"
	if existed {
		action = "Updated"
	}
	log.Printf("%s variable: %s = %s", action, name, value)

	if err := g.SaveToFile(); err != nil {
		log.Printf("ERROR: Failed to save variables after Set: %v", err)
		return err
	}
	return nil
}

// Get retrieves a variable's value.
// Returns the value and true if found, empty string and false otherwise.
func (g *Global) Get(name string) (string, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	val, ok := g.variables[name]
	return val, ok
}

// Unset removes a variable and saves to disk.
// Returns true if the variable existed and was removed, false otherwise.
func (g *Global) Unset(name string) (bool, error) {
	g.mu.Lock()
	existed := false
	if _, exists := g.variables[name]; exists {
		delete(g.variables, name)
		existed = true
	}
	g.mu.Unlock()

	if existed {
		log.Printf("Unset variable: %s", name)
		if err := g.SaveToFile(); err != nil {
			log.Printf("ERROR: Failed to save variables after Unset: %v", err)
			return true, err
		}
		return true, nil
	}
	log.Printf("Attempted to unset non-existent variable: %s", name)
	return false, nil
}

// All returns a copy of all variables, safe for iteration without holding the lock.
func (g *Global) All() map[string]string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make(map[string]string, len(g.variables))
	for k, v := range g.variables {
		result[k] = v
	}
	return result
}

// Count returns the number of stored variables.
func (g *Global) Count() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.variables)
}

// GetFilePath returns the path to the variables.json file.
func (g *Global) GetFilePath() string {
	return g.filePath
}

// =============================================================================
// Command Integration
// =============================================================================

// ApplyToCommand replaces {{var}} placeholders with stored values.
// Returns the modified command and a list of applied variable names.
func (g *Global) ApplyToCommand(command string) (string, []string) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var applied []string
	result := command

	for name, value := range g.variables {
		placeholder := "{{" + name + "}}"
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, value)
			applied = append(applied, name)
		}
	}

	return result, applied
}
