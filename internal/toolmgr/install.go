// Package toolmgr handles tool availability checking and installation.
package toolmgr

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hamedhdd/arsenal-ng/internal/model"
)

// DetectPackageManager returns the package manager command for the current system.
func DetectPackageManager() string {
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			return "brew"
		}
		return ""
	}

	// Linux package managers
	managers := []string{"apt", "dnf", "pacman", "yum"}
	for _, mgr := range managers {
		if _, err := exec.LookPath(mgr); err == nil {
			return mgr
		}
	}

	return ""
}

// GetInstallCommand returns the installation command for a tool.
// Returns empty string if no installation method is available.
func GetInstallCommand(install *model.InstallCmd) string {
	if install == nil {
		return ""
	}

	// Try language-specific package managers first (usually most reliable)
	if install.Go != "" {
		if _, err := exec.LookPath("go"); err == nil {
			return "go install " + install.Go
		}
	}

	if install.Pipx != "" {
		if _, err := exec.LookPath("pipx"); err == nil {
			return "pipx install " + install.Pipx
		}
	}

	if install.Pip != "" {
		if _, err := exec.LookPath("pip"); err == nil {
			return "pip install " + install.Pip
		}
		if _, err := exec.LookPath("pip3"); err == nil {
			return "pip3 install " + install.Pip
		}
	}

	if install.Cargo != "" {
		if _, err := exec.LookPath("cargo"); err == nil {
			return "cargo install " + install.Cargo
		}
	}

	// System package managers
	pkgMgr := DetectPackageManager()
	switch pkgMgr {
	case "brew":
		if install.Brew != "" {
			return "brew install " + install.Brew
		}
	case "apt":
		if install.Apt != "" {
			return "sudo apt install -y " + install.Apt
		}
	case "dnf", "yum":
		if install.Dnf != "" {
			return "sudo dnf install -y " + install.Dnf
		}
	case "pacman":
		if install.Pacman != "" {
			return "sudo pacman -S --noconfirm " + install.Pacman
		}
	}

	// Git clone as fallback
	if install.Git != "" {
		return install.Git
	}

	// Manual instructions
	if install.Manual != "" {
		return install.Manual
	}

	return ""
}

// FormatInstallPrompt returns a user-friendly prompt for tool installation.
func FormatInstallPrompt(tool string, install *model.InstallCmd) string {
	cmd := GetInstallCommand(install)
	if cmd == "" {
		return fmt.Sprintf("Tool '%s' is not installed and no installation method is available.\nPlease install it manually.", tool)
	}

	// Check if it's a multi-step process (git clone)
	if strings.Contains(cmd, "git clone") || strings.Contains(cmd, "&&") {
		return fmt.Sprintf("Tool '%s' is not installed.\n\nInstallation requires multiple steps:\n%s\n\nWould you like to proceed?", tool, cmd)
	}

	return fmt.Sprintf("Tool '%s' is not installed.\n\nInstall with: %s\n\nProceed with installation?", tool, cmd)
}
