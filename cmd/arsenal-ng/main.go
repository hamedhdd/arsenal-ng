// Package main is the entry point for arsenal-ng, a modern pentest command launcher.
//
// This program loads cheat files, initializes the TUI, and outputs the selected
// command to the terminal. Inspired by https://github.com/Orange-Cyberdefense/arsenal
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/atotto/clipboard"

	"github.com/hamedhdd/arsenal-ng/internal/config"
	"github.com/hamedhdd/arsenal-ng/internal/loader"
	"github.com/hamedhdd/arsenal-ng/internal/model"
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
	profile := flag.String("profile", "", "Filter tools by a specific tag (e.g., 'web', 'recon')")
	flag.Parse()

	if err := setupLogging(); err != nil {
		// Non-fatal, continue execution
		fmt.Fprintf(os.Stderr, "Warning: failed to setup logging: %v\n", err)
	}

	cheats, err := loadAndFilterCheats(*profile)
	if err != nil {
		return err
	}

	appModel, err := runTUI(cheats)
	if err != nil {
		return err
	}

	if appModel.Cancelled || appModel.FinalCommand == "" {
		log.Printf("No command to execute")
		return nil
	}

	if err := checkAndInstallTool(appModel.SelectedCheat); err != nil {
		return err
	}

	return handleCommandAction(appModel.FinalCommand)
}

// setupLogging initializes application logging to a file.
// Returns an error if logging cannot be set up, but the error is non-fatal.
func setupLogging() error {
	logPath, err := config.GetLogPath()
	if err != nil {
		return fmt.Errorf("failed to get log path: %w", err)
	}

	_, err = tea.LogToFile(logPath, config.AppName)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	// Note: The log file is intentionally not closed as it needs to persist
	// for the lifetime of the application. It will be closed when the program exits.
	log.Printf("Application started, logging to: %s", logPath)
	return nil
}

// loadAndFilterCheats loads all cheat files and optionally filters by profile tag.
func loadAndFilterCheats(profile string) ([]*model.Cheat, error) {
	log.Printf("Loading cheat files...")
	cheats, err := loader.Load()
	if err != nil {
		log.Printf("ERROR: Failed to load cheats: %v", err)
		return nil, fmt.Errorf("failed to load cheats: %w", err)
	}

	if profile != "" {
		log.Printf("Filtering cheats by profile: %s", profile)
		cheats = filterCheatsByTag(cheats, profile)
	}

	log.Printf("Loaded %d cheat(s) successfully", len(cheats))

	if len(cheats) == 0 {
		log.Printf("ERROR: No cheats found")
		return nil, fmt.Errorf("no cheats found")
	}

	return cheats, nil
}

// filterCheatsByTag returns only the cheats that have the specified tag.
func filterCheatsByTag(cheats []*model.Cheat, tag string) []*model.Cheat {
	filtered := make([]*model.Cheat, 0)
	for _, c := range cheats {
		for _, t := range c.Tags {
			if t == tag {
				filtered = append(filtered, c)
				break
			}
		}
	}
	return filtered
}

// runTUI starts the terminal UI and returns the user's selection.
func runTUI(cheats []*model.Cheat) (ui.App, error) {
	log.Printf("Starting TUI...")
	app := ui.New(cheats)
	program := tea.NewProgram(app, tea.WithAltScreen())

	result, err := program.Run()
	if err != nil {
		log.Printf("ERROR: TUI error: %v", err)
		return ui.App{}, fmt.Errorf("TUI error: %w", err)
	}

	appModel := result.(ui.App)

	if appModel.Cancelled {
		log.Printf("Application cancelled by user")
	} else if appModel.FinalCommand == "" {
		log.Printf("No command selected")
	}

	return appModel, nil
}

// checkAndInstallTool verifies if the required tool is installed.
// If not installed and installation info is available, prompts the user to install it.
func checkAndInstallTool(cheat *model.Cheat) error {
	if cheat == nil {
		return nil
	}

	toolBinary := toolmgr.GetToolBinary(cheat.Tool)
	if toolmgr.IsInstalled(toolBinary) {
		return nil
	}

	log.Printf("Tool '%s' is not installed", toolBinary)

	if cheat.InstallInfo == nil {
		fmt.Printf("\n⚠️  Tool '%s' is not installed. Please install it before using this command.\n", toolBinary)
		return nil
	}

	installCmd := toolmgr.GetInstallCommand(cheat.InstallInfo)
	if installCmd == "" {
		fmt.Printf("\n⚠️  Tool '%s' is not installed and no installation method is available.\n", toolBinary)
		fmt.Println("Please install it manually before using this command.")
		return nil
	}

	return offerToolInstallation(toolBinary, installCmd)
}

// offerToolInstallation prompts the user to install a missing tool.
func offerToolInstallation(toolBinary, installCmd string) error {
	fmt.Printf("\n⚠️  Tool '%s' is not installed.\n\n", toolBinary)
	fmt.Printf("Install with: %s\n\n", installCmd)
	fmt.Print("Would you like to install it now? [y/N]: ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		log.Printf("User declined installation")
		fmt.Println("\nℹ️  Skipping installation. Install the tool manually to use this command.")
		return nil
	}

	log.Printf("User confirmed installation of %s", toolBinary)
	fmt.Printf("\n🔧 Installing %s...\n\n", toolBinary)

	installExec := exec.Command("sh", "-c", installCmd)
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
	return nil
}

// handleCommandAction prompts the user for how to handle the final command
// (execute, inject to terminal, copy to clipboard, or cancel).
func handleCommandAction(command string) error {
	fmt.Printf("\n🚀 Command ready: %s\n\n", command)
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
		return executeCommand(command)
	case "2":
		injectToTerminal(command)
	case "3":
		return copyToClipboard(command)
	case "4", "q", "quit":
		log.Printf("User cancelled at action prompt")
		fmt.Println("Cancelled.")
	default:
		fmt.Println("Invalid choice. Cancelled.")
	}

	log.Printf("Application completed successfully")
	return nil
}

// executeCommand runs the command directly in a shell.
func executeCommand(command string) error {
	log.Printf("Executing command directly: %s", command)
	fmt.Printf("\nExecuting...\n\n")

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n❌ Execution failed: %v\n", err)
		return err
	}

	return nil
}

// injectToTerminal injects the command into the terminal input buffer using TIOCSTI.
func injectToTerminal(command string) {
	log.Printf("Outputting command to terminal: %s", command)
	output.ToTerminal(command)
}

// copyToClipboard copies the command to the system clipboard.
func copyToClipboard(command string) error {
	if err := clipboard.WriteAll(command); err != nil {
		fmt.Printf("\n❌ Failed to copy to clipboard: %v\n", err)
		log.Printf("Clipboard error: %v", err)
		return err
	}

	fmt.Printf("\n📋 Command copied to clipboard!\n")
	log.Printf("Command copied to clipboard")
	return nil
}

