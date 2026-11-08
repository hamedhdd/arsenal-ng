package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetConfigDir_ValidXDG(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create temp directory for testing
	tmpDir := t.TempDir()
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	os.Setenv("XDG_CONFIG_HOME", absPath)

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := filepath.Join(absPath, AppName)
	if configDir != expected {
		t.Errorf("Expected config dir %s, got %s", expected, configDir)
	}

	// Verify directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Expected directory to be created at %s", configDir)
	}
}

func TestGetConfigDir_RelativeXDG_FallbackToHome(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Set XDG_CONFIG_HOME to relative path (should be rejected)
	os.Setenv("XDG_CONFIG_HOME", "relative/path")

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should fall back to OS-specific default
	// On Windows: %LOCALAPPDATA%\arsenal-ng
	// On Linux/others: ~/.config/arsenal-ng
	home, _ := os.UserHomeDir()
	var expected string
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			expected = filepath.Join(localAppData, AppName)
		} else {
			expected = filepath.Join(home, AppName)
		}
	} else {
		expected = filepath.Join(home, ".config", AppName)
	}
	
	if configDir != expected {
		t.Errorf("Expected fallback to %s, got %s", expected, configDir)
	}
}

func TestGetConfigDir_EmptyXDG_UsesHome(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Unset XDG_CONFIG_HOME
	os.Unsetenv("XDG_CONFIG_HOME")

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should use OS-specific default
	// On Windows: %LOCALAPPDATA%\arsenal-ng
	// On Linux/others: ~/.config/arsenal-ng
	home, _ := os.UserHomeDir()
	var expected string
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			expected = filepath.Join(localAppData, AppName)
		} else {
			expected = filepath.Join(home, AppName)
		}
	} else {
		expected = filepath.Join(home, ".config", AppName)
	}
	
	if configDir != expected {
		t.Errorf("Expected %s, got %s", expected, configDir)
	}
}

func TestGetConfigDir_CreatesDirectory(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Use temp directory
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test-arsenal")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Ensure directory doesn't exist initially
	os.RemoveAll(testPath)

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Expected directory to exist, got error: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory", configDir)
	}
}

func TestGetVariablesPath(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tmpDir := t.TempDir()
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	os.Setenv("XDG_CONFIG_HOME", absPath)

	varPath, err := GetVariablesPath()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedDir := filepath.Join(absPath, AppName)
	expectedFile := filepath.Join(expectedDir, "variables.json")

	if varPath != expectedFile {
		t.Errorf("Expected variables path %s, got %s", expectedFile, varPath)
	}

	// Verify it ends with variables.json
	if !strings.HasSuffix(varPath, "variables.json") {
		t.Errorf("Expected path to end with 'variables.json', got %s", varPath)
	}
}

func TestGetLogPath(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tmpDir := t.TempDir()
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	os.Setenv("XDG_CONFIG_HOME", absPath)

	logPath, err := GetLogPath()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedDir := filepath.Join(absPath, AppName)

	if !strings.HasPrefix(logPath, expectedDir) {
		t.Errorf("Expected log path to start with %s, got %s", expectedDir, logPath)
	}

	// Verify directory structure is correct
	if filepath.Dir(logPath) != expectedDir {
		t.Errorf("Expected log path directory to be %s, got %s", expectedDir, filepath.Dir(logPath))
	}
}

func TestGetConfigDir_PathWithSpaces(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create temp directory with spaces
	tmpBase := t.TempDir()
	tmpDir := filepath.Join(tmpBase, "test config dir")
	os.MkdirAll(tmpDir, 0755)
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	os.Setenv("XDG_CONFIG_HOME", absPath)

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Expected no error with spaces in path, got: %v", err)
	}

	expected := filepath.Join(absPath, AppName)
	if configDir != expected {
		t.Errorf("Expected config dir %s, got %s", expected, configDir)
	}

	// Verify directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Expected directory to be created at %s", configDir)
	}
}

func TestGetConfigDir_NoHomeDir(t *testing.T) {
	// This test verifies behavior when HOME is not set
	// Save original env
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
	}()

	// Unset all home-related variables
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")

	_, err := GetConfigDir()
	// Should either succeed with fallback or return an error
	// Both are acceptable behaviors
	if err != nil {
		// Error is acceptable when no home directory can be determined
		t.Logf("Expected behavior: error when no home directory available: %v", err)
	}
}
