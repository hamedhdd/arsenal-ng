package loader

import (
	"testing"

	"github.com/hamedhdd/arsenal-ng/internal/model"
)

func TestSearch_EmptyQuery(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "nmap", Title: "nmap scan", Tags: []string{"scan"}},
		{Tool: "ffuf", Title: "ffuf dir", Tags: []string{"web"}},
	}

	result := Search(cheats, "")
	if len(result) != 2 {
		t.Errorf("expected all cheats returned for empty query, got %d", len(result))
	}
}

func TestSearch_SingleTerm(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "nmap", Title: "port scan", Tags: []string{"network"}},
		{Tool: "ffuf", Title: "dir scan", Tags: []string{"web"}},
	}

	result := Search(cheats, "nmap")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Tool != "nmap" {
		t.Errorf("expected nmap, got %s", result[0].Tool)
	}
}

func TestSearch_MultipleTerms(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "nmap", Title: "port scan", Tags: []string{"network"}},
		{Tool: "nmap", Title: "stealth scan", Tags: []string{"network"}},
		{Tool: "ffuf", Title: "dir scan", Tags: []string{"web"}},
	}

	// Both terms must match
	result := Search(cheats, "nmap port")
	if len(result) != 1 {
		t.Fatalf("expected 1 result matching both terms, got %d", len(result))
	}
	if result[0].Title != "port scan" {
		t.Errorf("unexpected result: %s", result[0].Title)
	}
}

func TestSearch_SearchAcrossFields(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "nmap", Title: "scan", Tags: []string{"recon"}, Command: "nmap -sV {{ip}}"},
	}

	// Match on tool
	if len(Search(cheats, "nmap")) != 1 {
		t.Error("failed to match on tool")
	}
	// Match on tag
	if len(Search(cheats, "recon")) != 1 {
		t.Error("failed to match on tag")
	}
	// Match on command
	if len(Search(cheats, "-sV")) != 1 {
		t.Error("failed to match on command")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "Nmap", Title: "Port Scan", Tags: []string{"Network"}},
	}

	result := Search(cheats, "nmap")
	if len(result) != 1 {
		t.Error("search should be case insensitive")
	}
}

func TestSearch_NoMatch(t *testing.T) {
	cheats := []*model.Cheat{
		{Tool: "nmap", Title: "scan"},
		{Tool: "ffuf", Title: "dir"},
	}

	result := Search(cheats, "doesnotexist")
	if len(result) != 0 {
		t.Errorf("expected 0 results for non-matching query, got %d", len(result))
	}
}
