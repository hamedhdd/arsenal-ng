package toolmgr

import (
	"testing"
)

func TestIsInstalled_CommonTools(t *testing.T) {
	tests := []struct {
		name     string
		tool     string
		expected bool
	}{
		{
			name:     "go should be installed (used to build this test)",
			tool:     "go",
			expected: true, // go must be available to run these tests
		},
		{
			name:     "nonexistent tool should not be installed",
			tool:     "this-tool-definitely-does-not-exist-12345",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInstalled(tt.tool)
			if result != tt.expected {
				t.Errorf("IsInstalled(%s) = %v, expected %v", tt.tool, result, tt.expected)
			}
		})
	}
}

func TestGetToolBinary_Simple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple tool name",
			input:    "nmap",
			expected: "nmap",
		},
		{
			name:     "hyphenated tool name",
			input:    "impacket-psexec",
			expected: "impacket-psexec",
		},
		{
			name:     "uppercase should be lowercased",
			input:    "NMAP",
			expected: "nmap",
		},
		{
			name:     "mixed case",
			input:    "NmAp",
			expected: "nmap",
		},
		{
			name:     "whitespace should be trimmed",
			input:    "  nmap  ",
			expected: "nmap",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "tool with numbers",
			input:    "python3",
			expected: "python3",
		},
		{
			name:     "tool with underscore",
			input:    "my_tool",
			expected: "my_tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetToolBinary(tt.input)
			if result != tt.expected {
				t.Errorf("GetToolBinary(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsInstalled_EmptyString(t *testing.T) {
	// Empty string should return false
	result := IsInstalled("")
	if result {
		t.Errorf("IsInstalled(\"\") should return false, got true")
	}
}

func TestIsInstalled_WithPath(t *testing.T) {
	// Tool names with path separators should be handled correctly
	result := IsInstalled("/usr/bin/this-does-not-exist")
	if result {
		t.Errorf("IsInstalled with full path to non-existent tool should return false")
	}
}
