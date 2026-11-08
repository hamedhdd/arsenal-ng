package ui

import (
	"strings"
	"testing"
)

func TestGetTagColor_Consistency(t *testing.T) {
	// Same tag should always return the same color
	tag := "web"
	color1 := getTagColor(tag)
	color2 := getTagColor(tag)

	if color1 != color2 {
		t.Errorf("getTagColor(%q) returned different colors: %s vs %s", tag, color1, color2)
	}
}

func TestGetTagColor_DifferentTags(t *testing.T) {
	// Different tags should generally return different colors
	// (though hash collisions are possible)
	tags := []string{"web", "recon", "exploit", "network", "database"}
	colors := make(map[string]string)

	for _, tag := range tags {
		color := getTagColor(tag)
		colors[tag] = color

		// Verify color is in hex format
		if !strings.HasPrefix(color, "#") {
			t.Errorf("getTagColor(%q) = %q, expected hex color starting with #", tag, color)
		}

		if len(color) != 7 {
			t.Errorf("getTagColor(%q) = %q, expected 7 characters (#RRGGBB)", tag, color)
		}
	}

	// Count unique colors
	uniqueColors := make(map[string]bool)
	for _, color := range colors {
		uniqueColors[color] = true
	}

	// With 5 tags, we should have at least 3 unique colors (allowing some collisions)
	if len(uniqueColors) < 3 {
		t.Errorf("Expected at least 3 unique colors for 5 different tags, got %d", len(uniqueColors))
	}
}

func TestGetTagColor_EmptyString(t *testing.T) {
	color := getTagColor("")

	// Should return a valid color even for empty string
	if !strings.HasPrefix(color, "#") {
		t.Errorf("getTagColor(\"\") = %q, expected hex color", color)
	}
}

func TestGetTagColor_CaseSensitivity(t *testing.T) {
	// Test if the function is case-sensitive
	color1 := getTagColor("web")
	color2 := getTagColor("WEB")
	color3 := getTagColor("Web")

	// These should be different (case-sensitive)
	// If they're the same, the function might be normalizing case
	allColors := []string{color1, color2, color3}
	for i, c := range allColors {
		if !strings.HasPrefix(c, "#") || len(c) != 7 {
			t.Errorf("Color %d is invalid: %s", i, c)
		}
	}
}

func TestGetToolColor_Consistency(t *testing.T) {
	// Same tool should always return the same color
	tool := "nmap"
	color1 := getToolColor(tool)
	color2 := getToolColor(tool)

	if color1 != color2 {
		t.Errorf("getToolColor(%q) returned different colors: %s vs %s", tool, color1, color2)
	}
}

func TestGetToolColor_Format(t *testing.T) {
	tools := []string{"nmap", "ffuf", "gobuster", "nuclei"}

	for _, tool := range tools {
		color := getToolColor(tool)

		if !strings.HasPrefix(color, "#") {
			t.Errorf("getToolColor(%q) = %q, expected hex color starting with #", tool, color)
		}

		if len(color) != 7 {
			t.Errorf("getToolColor(%q) = %q, expected 7 characters (#RRGGBB)", tool, color)
		}

		// Verify it's a valid hex color
		for i, c := range color[1:] {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
				t.Errorf("getToolColor(%q) = %q, character %d (%c) is not valid hex", tool, color, i+1, c)
			}
		}
	}
}

func TestRenderColoredTags_EmptySlice(t *testing.T) {
	tags := []string{}
	result := renderColoredTags(tags)

	if result != "" {
		t.Errorf("renderColoredTags([]) = %q, expected empty string", result)
	}
}

func TestRenderColoredTags_SingleTag(t *testing.T) {
	tags := []string{"web"}
	result := renderColoredTags(tags)

	if result == "" {
		t.Errorf("renderColoredTags([\"web\"]) returned empty string")
	}

	// Should contain the tag name
	if !strings.Contains(result, "web") {
		t.Errorf("renderColoredTags([\"web\"]) = %q, should contain 'web'", result)
	}
}

func TestRenderColoredTags_MultipleTags(t *testing.T) {
	tags := []string{"web", "recon", "exploit"}
	result := renderColoredTags(tags)

	if result == "" {
		t.Errorf("renderColoredTags(%v) returned empty string", tags)
	}

	// Should contain all tag names
	for _, tag := range tags {
		if !strings.Contains(result, tag) {
			t.Errorf("renderColoredTags(%v) = %q, should contain '%s'", tags, result, tag)
		}
	}
}

func TestRenderColoredTags_ManyTags(t *testing.T) {
	// Test with a large number of tags
	tags := []string{"web", "recon", "exploit", "network", "database", "enum", "scan", "post-exploit"}
	result := renderColoredTags(tags)

	if result == "" {
		t.Errorf("renderColoredTags with many tags returned empty string")
	}

	// Should contain at least some of the tags
	foundCount := 0
	for _, tag := range tags {
		if strings.Contains(result, tag) {
			foundCount++
		}
	}

	if foundCount == 0 {
		t.Errorf("renderColoredTags should contain at least some tags, found none")
	}
}

func TestRenderColoredTags_NilSlice(t *testing.T) {
	var tags []string = nil
	result := renderColoredTags(tags)

	// Should handle nil slice gracefully
	if result != "" {
		t.Errorf("renderColoredTags(nil) = %q, expected empty string", result)
	}
}

func TestRenderColoredTags_DuplicateTags(t *testing.T) {
	tags := []string{"web", "web", "recon"}
	result := renderColoredTags(tags)

	// Should handle duplicates (exact behavior may vary)
	if result == "" {
		t.Errorf("renderColoredTags with duplicates returned empty string")
	}
}

func TestRenderColoredTags_EmptyStringTag(t *testing.T) {
	tags := []string{"web", "", "recon"}
	result := renderColoredTags(tags)

	// Should handle empty strings in the slice
	if !strings.Contains(result, "web") || !strings.Contains(result, "recon") {
		t.Errorf("renderColoredTags should contain non-empty tags")
	}
}
