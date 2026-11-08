//go:build windows

// Package output provides terminal output functionality.
//
// On Windows, TIOCSTI is not available. Commands are printed to stdout
// instead of being injected into the terminal input buffer.
package output

import (
	"fmt"
	"log"
)

// ToTerminal prints the command to stdout on Windows.
// TIOCSTI terminal injection is not supported on Windows.
// The command is printed on a new line for the user to copy and execute.
func ToTerminal(command string) {
	if len(command) == 0 {
		log.Printf("WARNING: Attempted to output empty command to terminal")
		return
	}

	log.Printf("TIOCSTI not available on Windows — printing command to stdout")
	fmt.Printf("\n%s\n", command)
}
