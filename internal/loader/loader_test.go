package loader

import (
	"testing"
)

// TestLoad_AllCheatFilesValid ensures every embedded cheat file parses
// without error and produces at least one valid Cheat entry.
// This acts as a CI gate against broken YAML contributions.
func TestLoad_AllCheatFilesValid(t *testing.T) {
	cheats, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if len(cheats) == 0 {
		t.Fatal("Load() returned 0 cheats — expected embedded cheat files to be loaded")
	}
	t.Logf("Loaded %d cheat(s) from embedded files", len(cheats))
}

// TestLoad_AllCheatsHaveRequiredFields ensures every cheat has the
// minimum fields needed to be useful: a tool name and a command.
func TestLoad_AllCheatsHaveRequiredFields(t *testing.T) {
	cheats, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	for _, c := range cheats {
		if c.Tool == "" {
			t.Errorf("cheat from file '%s' has empty tool name (title: '%s')", c.Filename, c.Title)
		}
		if c.Command == "" {
			t.Errorf("cheat '%s' in tool '%s' has empty command", c.Title, c.Tool)
		}
	}
}

// TestLoad_AllCheatsHaveTitles ensures every cheat action has a display title.
func TestLoad_AllCheatsHaveTitles(t *testing.T) {
	cheats, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	for _, c := range cheats {
		if c.Title == "" {
			t.Errorf("cheat in tool '%s' (file: '%s') has no title", c.Tool, c.Filename)
		}
	}
}

// TestLoad_Idempotent ensures calling Load() twice returns the same count.
func TestLoad_Idempotent(t *testing.T) {
	cheats1, _ := Load()
	cheats2, _ := Load()
	if len(cheats1) != len(cheats2) {
		t.Errorf("Load() is not idempotent: first call returned %d, second returned %d", len(cheats1), len(cheats2))
	}
}
