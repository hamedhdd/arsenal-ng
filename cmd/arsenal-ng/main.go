// Package main is the entry point for arsenal-ng, a modern pentest command launcher.
//
// This program loads cheat files, initializes the TUI, and outputs the selected
// command to the terminal. Inspired by https://github.com/Orange-Cyberdefense/arsenal
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/atotto/clipboard"

	"github.com/hamedhdd/arsenal-ng/internal/config"
	"github.com/hamedhdd/arsenal-ng/internal/loader"
	"github.com/hamedhdd/arsenal-ng/internal/output"
	"github.com/hamedhdd/arsenal-ng/internal/toolmgr"
	"github.com/hamedhdd/arsenal-ng/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Setup logging to file (same directory as variables.json)
	logPath, err := config.GetLogPath()
	if err != nil {
		// Log error but don't fail - logging is optional
		fmt.Fprintf(os.Stderr, "Warning: failed to get log path: %v\n", err)
	} else {
		logFile, err := tea.LogToFile(logPath, config.AppName)
		if err != nil {
			// Log error but don't fail - logging is optional
			fmt.Fprintf(os.Stderr, "Warning: failed to setup logging: %v\n", err)
		} else {
			defer logFile.Close()
			log.Printf("Application started, logging to: %s", logPath)
		}
	}

	// Load cheat files
	log.Printf("Loading cheat files...")
	cheats, err := loader.Load()
	if err != nil {
		log.Printf("ERROR: Failed to load cheats: %v", err)
		return fmt.Errorf("failed to load cheats: %w", err)
	}

	log.Printf("Loaded %d cheat(s) successfully", len(cheats))

	if len(cheats) == 0 {
		log.Printf("ERROR: No cheats found")
		return fmt.Errorf("no cheats found")
	}

	// Run TUI
	log.Printf("Starting TUI...")
	app := ui.New(cheats)
	program := tea.NewProgram(app, tea.WithAltScreen())

	result, err := program.Run()
	if err != nil {
		log.Printf("ERROR: TUI error: %v", err)
		return fmt.Errorf("TUI error: %w", err)
	}

	// Handle result
	model := result.(ui.App)

	if model.Cancelled {
		log.Printf("Application cancelled by user")
		return nil
	}

	if model.FinalCommand == "" {
		log.Printf("No command selected")
		return nil
	}

	// Check if tool is installed (if we have the selected cheat info)
	if model.SelectedCheat != nil {
		toolBinary := toolmgr.GetToolBinary(model.SelectedCheat.Tool)
		if !toolmgr.IsInstalled(toolBinary) {
			log.Printf("Tool '%s' is not installed", toolBinary)

			// Offer installation if we have install info
			if model.SelectedCheat.InstallInfo != nil {
				installCmd := toolmgr.GetInstallCommand(model.SelectedCheat.InstallInfo)
				if installCmd != "" {
					fmt.Printf("\n⚠️  Tool '%s' is not installed.\n\n", toolBinary)
					fmt.Printf("Install with: %s\n\n", installCmd)
					fmt.Print("Would you like to install it now? [y/N]: ")

					var response string
					fmt.Scanln(&response)
					response = strings.ToLower(strings.TrimSpace(response))

					if response == "y" || response == "yes" {
						log.Printf("User confirmed installation of %s", toolBinary)
						fmt.Printf("\n🔧 Installing %s...\n\n", toolBinary)

						// Execute installation command
						cmdParts := strings.Fields(installCmd)
						if len(cmdParts) > 0 {
							installExec := exec.Command(cmdParts[0], cmdParts[1:]...)
							installExec.Stdout = os.Stdout
							installExec.Stderr = os.Stderr
							installExec.Stdin = os.Stdin

							if err := installExec.Run(); err != nil {
								log.Printf("ERROR: Installation failed: %v", err)
								fmt.Printf("\n❌ Installation failed: %v\n", err)
								fmt.Printf("Please install %s manually.\n", toolBinary)
								return nil
							}

							fmt.Printf("\n✅ %s installed successfully!\n\n", toolBinary)
							log.Printf("%s installed successfully", toolBinary)
						}
					} else {
						log.Printf("User declined installation")
						fmt.Println("\nℹ️  Skipping installation. Install the tool manually to use this command.")
						return nil
					}
				} else {
					fmt.Printf("\n⚠️  Tool '%s' is not installed and no installation method is available.\n", toolBinary)
					fmt.Println("Please install it manually before using this command.")
					return nil
				}
			} else {
				fmt.Printf("\n⚠️  Tool '%s' is not installed. Please install it before using this command.\n", toolBinary)
				return nil
			}
		}
	}

	// Action prompt
	fmt.Printf("\n🚀 Command ready: %s\n\n", model.FinalCommand)
	fmt.Println("How would you like to proceed?")
	fmt.Println("  [1] Execute directly")
	fmt.Println("  [2] Inject into terminal (TIOCSTI)")
	fmt.Println("  [3] Copy to clipboard")
	fmt.Println("  [4] Cancel")
	fmt.Print("\nChoice [1]: ")

	var action string
	fmt.Scanln(&action)
	action = strings.TrimSpace(action)

	switch action {
	case "1", "":
		log.Printf("Executing command directly: %s", model.FinalCommand)
		fmt.Printf("\nExecuting...\n\n")
		cmd := exec.Command("sh", "-c", model.FinalCommand)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("\n❌ Execution failed: %v\n", err)
		}
	case "2":
		log.Printf("Outputting command to terminal: %s", model.FinalCommand)
		output.ToTerminal(model.FinalCommand)
	case "3":
		if err := clipboard.WriteAll(model.FinalCommand); err != nil {
			fmt.Printf("\n❌ Failed to copy to clipboard: %v\n", err)
			log.Printf("Clipboard error: %v", err)
		} else {
			fmt.Printf("\n📋 Command copied to clipboard!\n")
			log.Printf("Command copied to clipboard")
		}
	case "4", "q", "quit":
		log.Printf("User cancelled at action prompt")
		fmt.Println("Cancelled.")
	default:
		fmt.Println("Invalid choice. Cancelled.")
	}

	log.Printf("Application completed successfully")

	return nil
}

