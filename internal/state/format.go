// Package state manages persistent and session-scoped state for the application.
//
// This file provides display/formatting functions for the variables store.
package state

import (
	"fmt"
	"sort"
	"strings"
)

// FormatList returns a formatted string of all variables for display.
func (g *Global) FormatList() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if len(g.variables) == 0 {
		return "No variables set"
	}

	keys := make([]string, 0, len(g.variables))
	for k := range g.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Variables (%d):\n", len(g.variables)))
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  %s = %s\n", k, g.variables[k]))
	}
	return sb.String()
}
