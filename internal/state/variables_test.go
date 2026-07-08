package state

import (
	"os"
	"path/filepath"
	"testing"
)

// newTestGlobal creates a Global backed by a temp file for testing.
func newTestGlobal(t *testing.T) *Global {
	t.Helper()
	dir := t.TempDir()
	g := &Global{
		variables: make(map[string]string),
		filePath:  filepath.Join(dir, "variables.json"),
	}
	return g
}

// =============================================================================
// Set / Get / Unset
// =============================================================================

func TestSet_Get(t *testing.T) {
	g := newTestGlobal(t)
	if err := g.Set("ip", "10.10.10.10"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	val, ok := g.Get("ip")
	if !ok {
		t.Fatal("Get returned not found after Set")
	}
	if val != "10.10.10.10" {
		t.Errorf("expected '10.10.10.10', got '%s'", val)
	}
}

func TestSet_Overwrite(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "1.1.1.1")
	_ = g.Set("ip", "2.2.2.2")
	val, _ := g.Get("ip")
	if val != "2.2.2.2" {
		t.Errorf("expected overwritten value '2.2.2.2', got '%s'", val)
	}
}

func TestGet_Missing(t *testing.T) {
	g := newTestGlobal(t)
	_, ok := g.Get("nonexistent")
	if ok {
		t.Error("expected Get to return false for missing key")
	}
}

func TestUnset_Existing(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.10.10.10")
	existed, err := g.Unset("ip")
	if err != nil {
		t.Fatalf("Unset failed: %v", err)
	}
	if !existed {
		t.Error("expected existed=true for key that was set")
	}
	_, ok := g.Get("ip")
	if ok {
		t.Error("expected key to be removed after Unset")
	}
}

func TestUnset_Missing(t *testing.T) {
	g := newTestGlobal(t)
	existed, err := g.Unset("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if existed {
		t.Error("expected existed=false for key that never existed")
	}
}

func TestCount(t *testing.T) {
	g := newTestGlobal(t)
	if g.Count() != 0 {
		t.Error("expected 0 on new store")
	}
	_ = g.Set("a", "1")
	_ = g.Set("b", "2")
	if g.Count() != 2 {
		t.Errorf("expected 2, got %d", g.Count())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.0.0.1")
	all := g.All()
	// Mutating the copy must not affect the store
	all["ip"] = "evil"
	val, _ := g.Get("ip")
	if val != "10.0.0.1" {
		t.Error("All() should return an isolated copy, not a reference")
	}
}

// =============================================================================
// ApplyToCommand
// =============================================================================

func TestApplyToCommand_Replaces(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.10.10.10")
	result, applied := g.ApplyToCommand("nmap -sV {{ip}}")
	if result != "nmap -sV 10.10.10.10" {
		t.Errorf("unexpected result: %s", result)
	}
	if len(applied) != 1 || applied[0] != "ip" {
		t.Errorf("unexpected applied list: %v", applied)
	}
}

func TestApplyToCommand_NoMatch(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.10.10.10")
	result, applied := g.ApplyToCommand("nmap -sV {{target}}")
	if result != "nmap -sV {{target}}" {
		t.Errorf("command should be unchanged: %s", result)
	}
	if len(applied) != 0 {
		t.Errorf("expected no applied vars, got %v", applied)
	}
}

func TestApplyToCommand_Multiple(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.0.0.1")
	_ = g.Set("port", "445")
	result, applied := g.ApplyToCommand("crackmapexec smb {{ip}} -p {{port}}")
	expected := "crackmapexec smb 10.0.0.1 -p 445"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
	if len(applied) != 2 {
		t.Errorf("expected 2 applied vars, got %d", len(applied))
	}
}

// =============================================================================
// Persistence (SaveToFile / LoadFromFile)
// =============================================================================

func TestPersistence_RoundTrip(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.0.0.1")
	_ = g.Set("domain", "corp.local")

	// Load into a fresh instance using the same file
	g2 := &Global{
		variables: make(map[string]string),
		filePath:  g.filePath,
	}
	if err := g2.LoadFromFile(); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}
	if g2.Count() != 2 {
		t.Errorf("expected 2 variables after load, got %d", g2.Count())
	}
	val, _ := g2.Get("ip")
	if val != "10.0.0.1" {
		t.Errorf("expected '10.0.0.1', got '%s'", val)
	}
}

func TestPersistence_AtomicWrite(t *testing.T) {
	g := newTestGlobal(t)
	_ = g.Set("ip", "10.0.0.1")

	// .tmp file must not exist after a successful save
	tmpPath := g.filePath + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("temporary file should be cleaned up after atomic save")
	}
}

func TestPersistence_EmptyFile(t *testing.T) {
	g := newTestGlobal(t)
	// Write empty file
	_ = os.WriteFile(g.filePath, []byte{}, 0644)
	if err := g.LoadFromFile(); err != nil {
		t.Errorf("LoadFromFile should handle empty file gracefully: %v", err)
	}
	if g.Count() != 0 {
		t.Error("expected 0 variables from empty file")
	}
}

func TestPersistence_NoPersistence(t *testing.T) {
	// filePath="" means persistence disabled — no error expected
	g := &Global{variables: make(map[string]string)}
	if err := g.SaveToFile(); err != nil {
		t.Errorf("SaveToFile with empty path should be a no-op: %v", err)
	}
	if err := g.LoadFromFile(); err != nil {
		t.Errorf("LoadFromFile with empty path should be a no-op: %v", err)
	}
}
