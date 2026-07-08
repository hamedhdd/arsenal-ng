// Package toolmgr handles tool availability checking and installation.
package toolmgr

import (
	"os/exec"
	"strings"
)

// IsInstalled checks if a tool is available in PATH.
func IsInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// GetToolBinary extracts the binary name from a tool name.
// Handles cases like "impacket-psexec" → "impacket-psexec"
// and "nmap" → "nmap".
func GetToolBinary(tool string) string {
	// For most cases, the tool name is the binary name
	// Special cases can be added here if needed
	return strings.ToLower(strings.TrimSpace(tool))
}
