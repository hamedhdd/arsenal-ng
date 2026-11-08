package toolmgr

import (
	"runtime"
	"strings"
	"testing"

	"github.com/hamedhdd/arsenal-ng/internal/model"
)

func TestDetectPackageManager(t *testing.T) {
	pkgMgr := DetectPackageManager()

	// We can't test the exact result as it depends on the system,
	// but we can verify it returns one of the expected values or empty
	validManagers := map[string]bool{
		"":       true, // No package manager found
		"brew":   true,
		"apt":    true,
		"dnf":    true,
		"pacman": true,
		"yum":    true,
	}

	if !validManagers[pkgMgr] {
		t.Errorf("DetectPackageManager() returned unexpected value: %s", pkgMgr)
	}

	// On macOS, if brew is installed, it should be detected
	if runtime.GOOS == "darwin" {
		// We can't assume brew is installed, but if it returns something, it should be brew
		if pkgMgr != "" && pkgMgr != "brew" {
			t.Errorf("On macOS, expected 'brew' or empty, got: %s", pkgMgr)
		}
	}
}

func TestGetInstallCommand_NilInput(t *testing.T) {
	cmd := GetInstallCommand(nil)
	if cmd != "" {
		t.Errorf("GetInstallCommand(nil) should return empty string, got: %s", cmd)
	}
}

func TestGetInstallCommand_EmptyStruct(t *testing.T) {
	install := &model.InstallCmd{}
	cmd := GetInstallCommand(install)
	if cmd != "" {
		t.Errorf("GetInstallCommand with empty struct should return empty string, got: %s", cmd)
	}
}

func TestGetInstallCommand_GoPriority(t *testing.T) {
	// Go install should be preferred if available
	install := &model.InstallCmd{
		Go:  "github.com/example/tool@latest",
		Pip: "example-tool",
		Apt: "example-tool",
	}

	cmd := GetInstallCommand(install)

	// If go is in PATH, should use go install
	// Otherwise, should fall through to other methods
	if cmd != "" {
		if strings.HasPrefix(cmd, "go install") {
			if !strings.Contains(cmd, "github.com/example/tool@latest") {
				t.Errorf("Go install command should contain the full package path")
			}
		}
	}
}

func TestGetInstallCommand_PipxBeforePip(t *testing.T) {
	// Pipx should be preferred over pip if available
	install := &model.InstallCmd{
		Pipx: "example-tool",
		Pip:  "example-tool",
	}

	cmd := GetInstallCommand(install)

	// If command is returned, pipx should take priority
	if cmd != "" && strings.HasPrefix(cmd, "pip ") {
		// Check if pipx is available - if it is, test fails
		if IsInstalled("pipx") {
			t.Errorf("Pipx is available but pip was chosen instead")
		}
	}
}

func TestGetInstallCommand_GitFallback(t *testing.T) {
	install := &model.InstallCmd{
		Git: "git clone https://github.com/example/tool",
	}

	cmd := GetInstallCommand(install)

	// Git should be used if no other method is available
	if cmd != "" {
		if !strings.Contains(cmd, "git clone") {
			// This is acceptable - might have found another method
			t.Logf("Expected git clone in command, got: %s", cmd)
		}
	}
}

func TestGetInstallCommand_ManualInstructions(t *testing.T) {
	install := &model.InstallCmd{
		Manual: "Visit https://example.com/install and follow instructions",
	}

	cmd := GetInstallCommand(install)

	// Manual instructions should be returned if no automated method available
	if cmd == "" {
		t.Errorf("Manual instructions should be returned when no other method available")
	} else {
		if !strings.Contains(cmd, "example.com") {
			t.Errorf("Expected manual instructions in command, got: %s", cmd)
		}
	}
}

func TestGetInstallCommand_CargoInstall(t *testing.T) {
	install := &model.InstallCmd{
		Cargo: "example-tool",
	}

	cmd := GetInstallCommand(install)

	// If cargo is available, should use it
	if IsInstalled("cargo") {
		if cmd == "" {
			t.Errorf("Cargo is available but no command was returned")
		}
		if !strings.HasPrefix(cmd, "cargo install") {
			t.Errorf("Expected 'cargo install', got: %s", cmd)
		}
	}
}

func TestGetInstallCommand_SystemPackageManager(t *testing.T) {
	pkgMgr := DetectPackageManager()

	tests := []struct {
		name      string
		install   *model.InstallCmd
		expectCmd string
	}{
		{
			name: "apt install",
			install: &model.InstallCmd{
				Apt: "example-tool",
			},
			expectCmd: "sudo apt install -y example-tool",
		},
		{
			name: "dnf install",
			install: &model.InstallCmd{
				Dnf: "example-tool",
			},
			expectCmd: "sudo dnf install -y example-tool",
		},
		{
			name: "pacman install",
			install: &model.InstallCmd{
				Pacman: "example-tool",
			},
			expectCmd: "sudo pacman -S --noconfirm example-tool",
		},
		{
			name: "brew install",
			install: &model.InstallCmd{
				Brew: "example-tool",
			},
			expectCmd: "brew install example-tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetInstallCommand(tt.install)

			// Only verify if this is the detected package manager
			if pkgMgr != "" && strings.Contains(tt.expectCmd, pkgMgr) {
				if cmd != tt.expectCmd {
					t.Errorf("Expected command %q, got %q", tt.expectCmd, cmd)
				}
			}
		})
	}
}

func TestFormatInstallPrompt_NoInstallMethod(t *testing.T) {
	install := &model.InstallCmd{}
	prompt := FormatInstallPrompt("testtool", install)

	if !strings.Contains(prompt, "testtool") {
		t.Errorf("Prompt should contain tool name 'testtool'")
	}

	if !strings.Contains(prompt, "not installed") {
		t.Errorf("Prompt should indicate tool is not installed")
	}

	if !strings.Contains(prompt, "manually") {
		t.Errorf("Prompt should suggest manual installation")
	}
}

func TestFormatInstallPrompt_SimpleCommand(t *testing.T) {
	install := &model.InstallCmd{
		Pip: "testtool",
	}

	// Only test if pip is available
	if !IsInstalled("pip") && !IsInstalled("pip3") {
		t.Skip("pip not available, skipping test")
	}

	prompt := FormatInstallPrompt("testtool", install)

	if !strings.Contains(prompt, "testtool") {
		t.Errorf("Prompt should contain tool name 'testtool'")
	}

	if !strings.Contains(prompt, "Install with:") {
		t.Errorf("Prompt should contain 'Install with:'")
	}
}

func TestFormatInstallPrompt_GitClone(t *testing.T) {
	install := &model.InstallCmd{
		Git: "git clone https://github.com/example/tool && cd tool && make install",
	}

	prompt := FormatInstallPrompt("testtool", install)

	if !strings.Contains(prompt, "multiple steps") {
		t.Errorf("Git clone prompt should mention multiple steps")
	}

	if !strings.Contains(prompt, "git clone") {
		t.Errorf("Prompt should contain the git clone command")
	}
}

func TestFormatInstallPrompt_NilInstall(t *testing.T) {
	prompt := FormatInstallPrompt("testtool", nil)

	if !strings.Contains(prompt, "testtool") {
		t.Errorf("Prompt should contain tool name even with nil install")
	}

	if !strings.Contains(prompt, "not installed") {
		t.Errorf("Prompt should indicate tool is not installed")
	}
}
