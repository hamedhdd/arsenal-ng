// Package state manages persistent and session-scoped state for the application.
//
// This file handles atomic file I/O for the variables store (LoadFromFile,
// SaveToFile). Uses temp file + rename for atomic writes.
package state

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LoadFromFile loads variables from the JSON file on disk.
// Returns os.ErrNotExist if file doesn't exist (expected on first run).
// Thread-safe.
func (g *Global) LoadFromFile() error {
	if g.filePath == "" {
		return nil // Persistence disabled
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	data, err := os.ReadFile(g.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil // Empty file is valid
	}

	var vars map[string]string
	if err := json.Unmarshal(data, &vars); err != nil {
		log.Printf("ERROR: Failed to parse variables.json: %v", err)
		return fmt.Errorf("failed to parse variables.json: %w", err)
	}

	g.variables = vars
	return nil
}

// SaveToFile saves variables to disk atomically.
// Thread-safe.
func (g *Global) SaveToFile() error {
	if g.filePath == "" {
		return nil // Persistence disabled
	}

	g.mu.RLock()
	varsCopy := make(map[string]string, len(g.variables))
	for k, v := range g.variables {
		varsCopy[k] = v
	}
	g.mu.RUnlock()

	data, err := json.MarshalIndent(varsCopy, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal variables: %v", err)
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	tmpPath := g.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		log.Printf("ERROR: Failed to write variables file: %v", err)
		return fmt.Errorf("failed to write variables file: %w", err)
	}

	if err := os.Rename(tmpPath, g.filePath); err != nil {
		_ = os.Remove(tmpPath)
		log.Printf("ERROR: Failed to rename variables file: %v", err)
		return fmt.Errorf("failed to rename variables file: %w", err)
	}

	log.Printf("Saved %d variable(s) to %s", len(varsCopy), g.filePath)
	return nil
}
