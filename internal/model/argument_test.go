package model

import (
	"testing"
)

// =============================================================================
// ParseArguments Tests
// =============================================================================

func TestParseArguments_Empty(t *testing.T) {
	args := ParseArguments("nmap -sV 192.168.1.1")
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestParseArguments_Single(t *testing.T) {
	args := ParseArguments("nmap -sV {{ip}}")
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	if args[0].Name != "ip" {
		t.Errorf("expected name 'ip', got '%s'", args[0].Name)
	}
	if args[0].DefaultValue != "" {
		t.Errorf("expected empty default, got '%s'", args[0].DefaultValue)
	}
}

func TestParseArguments_WithDefault(t *testing.T) {
	args := ParseArguments("nmap -p {{port|8080}} {{ip}}")
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0].Name != "port" {
		t.Errorf("expected name 'port', got '%s'", args[0].Name)
	}
	if args[0].DefaultValue != "8080" {
		t.Errorf("expected default '8080', got '%s'", args[0].DefaultValue)
	}
	if args[0].Value != "8080" {
		t.Errorf("expected value pre-filled with default '8080', got '%s'", args[0].Value)
	}
}

func TestParseArguments_Deduplication(t *testing.T) {
	// Same placeholder used twice — should return only one arg
	args := ParseArguments("nmap {{ip}} && curl {{ip}}")
	if len(args) != 1 {
		t.Errorf("expected 1 unique arg, got %d", len(args))
	}
}

func TestParseArguments_Multiple(t *testing.T) {
	args := ParseArguments("ffuf -u {{url}} -w {{wordlist}} -o {{output|scan.txt}}")
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	names := []string{args[0].Name, args[1].Name, args[2].Name}
	expected := []string{"url", "wordlist", "output"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("arg[%d]: expected name '%s', got '%s'", i, n, names[i])
		}
	}
}

func TestParseArguments_OrderPreserved(t *testing.T) {
	args := ParseArguments("{{c}} {{a}} {{b}}")
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[0].Name != "c" || args[1].Name != "a" || args[2].Name != "b" {
		t.Errorf("order not preserved: got %v, %v, %v", args[0].Name, args[1].Name, args[2].Name)
	}
}

// =============================================================================
// BuildCommand Tests
// =============================================================================

func TestBuildCommand_NoArgs(t *testing.T) {
	result := BuildCommand("nmap -sV 10.0.0.1", nil)
	if result != "nmap -sV 10.0.0.1" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestBuildCommand_SingleArg(t *testing.T) {
	args := ParseArguments("nmap -sV {{ip}}")
	args[0].Value = "10.10.10.10"
	result := BuildCommand("nmap -sV {{ip}}", args)
	if result != "nmap -sV 10.10.10.10" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestBuildCommand_WithDefault(t *testing.T) {
	args := ParseArguments("nmap -p {{port|8080}} {{ip}}")
	args[1].Value = "10.10.10.10"
	// port keeps its default
	result := BuildCommand("nmap -p {{port|8080}} {{ip}}", args)
	if result != "nmap -p 8080 10.10.10.10" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestBuildCommand_ReplacesAllOccurrences(t *testing.T) {
	args := ParseArguments("{{ip}} && {{ip}}")
	args[0].Value = "1.2.3.4"
	result := BuildCommand("{{ip}} && {{ip}}", args)
	if result != "1.2.3.4 && 1.2.3.4" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestBuildCommand_MultipleArgs(t *testing.T) {
	cmd := "ffuf -u {{url}} -w {{wordlist}}"
	args := ParseArguments(cmd)
	args[0].Value = "http://target.com/FUZZ"
	args[1].Value = "/usr/share/wordlists/common.txt"
	result := BuildCommand(cmd, args)
	expected := "ffuf -u http://target.com/FUZZ -w /usr/share/wordlists/common.txt"
	if result != expected {
		t.Errorf("unexpected result: %s", result)
	}
}

// =============================================================================
// HasEmptyArgs / GetEmptyArgs Tests
// =============================================================================

func TestHasEmptyArgs_AllFilled(t *testing.T) {
	args := []Argument{
		{Name: "ip", Value: "10.0.0.1"},
		{Name: "port", Value: "80"},
	}
	if HasEmptyArgs(args) {
		t.Error("expected no empty args")
	}
}

func TestHasEmptyArgs_OneEmpty(t *testing.T) {
	args := []Argument{
		{Name: "ip", Value: "10.0.0.1"},
		{Name: "port", Value: ""},
	}
	if !HasEmptyArgs(args) {
		t.Error("expected empty args to be detected")
	}
}

func TestGetEmptyArgs(t *testing.T) {
	args := []Argument{
		{Name: "ip", Value: "10.0.0.1"},
		{Name: "port", Value: ""},
		{Name: "user", Value: ""},
	}
	empty := GetEmptyArgs(args)
	if len(empty) != 2 {
		t.Fatalf("expected 2 empty args, got %d", len(empty))
	}
	if empty[0] != "port" || empty[1] != "user" {
		t.Errorf("unexpected empty args: %v", empty)
	}
}
