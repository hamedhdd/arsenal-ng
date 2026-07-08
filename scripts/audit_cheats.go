package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hamedhdd/arsenal-ng/internal/loader"
)

func main() {
	cheats, err := loader.Load()
	if err != nil {
		fmt.Printf("Error loading cheats: %v\n", err)
		os.Exit(1)
	}

	suspicious := 0
	totalInstalls := 0

	fmt.Println("Starting security audit of tool installation sources...")

	for _, cheat := range cheats {
		if cheat.InstallInfo == nil {
			continue
		}
		totalInstalls++
		
		info := cheat.InstallInfo
		
		// Check Git
		if info.Git != "" {
			if !strings.Contains(info.Git, "github.com") && !strings.Contains(info.Git, "gitlab.com") {
				fmt.Printf("🚨 [SUSPICIOUS] %s (Git): %s\n", cheat.Tool, info.Git)
				suspicious++
			}
		}
		
		// Check for malicious bash injections in package managers
		checkPackage := func(pkg, name string) {
			if pkg == "" { return }
			if strings.Contains(pkg, ";") || strings.Contains(pkg, "&&") || strings.Contains(pkg, "|") || strings.HasPrefix(pkg, "http") {
				fmt.Printf("🚨 [SUSPICIOUS] %s (%s): %s\n", cheat.Tool, name, pkg)
				suspicious++
			}
		}
		
		checkPackage(info.Pip, "Pip")
		checkPackage(info.Pipx, "Pipx")
		checkPackage(info.Brew, "Brew")
		checkPackage(info.Apt, "Apt")
		checkPackage(info.Dnf, "Dnf")
		checkPackage(info.Pacman, "Pacman")
		checkPackage(info.Go, "Go")
		checkPackage(info.Cargo, "Cargo")
	}

	fmt.Printf("\nAudit Complete! Checked %d tools with install instructions.\n", totalInstalls)
	if suspicious > 0 {
		fmt.Printf("❌ FAILED: Found %d suspicious entries. Please use valid domains and avoid chaining commands in package fields.\n", suspicious)
		os.Exit(1)
	}

	fmt.Println("✅ PASSED: All install sources appear legitimate.")
	os.Exit(0)
}
