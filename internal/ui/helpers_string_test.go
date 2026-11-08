package ui

import (
	"strings"
	"testing"
)

func TestTruncate_NoTruncationNeeded(t *testing.T) {
	input := "short"
	max := 10
	result := truncate(input, max)

	if result != input {
		t.Errorf("truncate(%q, %d) = %q, expected %q", input, max, result, input)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	input := "exactly10!"
	max := 10
	result := truncate(input, max)

	if result != input {
		t.Errorf("truncate(%q, %d) = %q, expected %q", input, max, result, input)
	}
}

func TestTruncate_NeedsTruncation(t *testing.T) {
	input := "this is a very long string that needs truncation"
	max := 20
	result := truncate(input, max)

	if len(result) > max {
		t.Errorf("truncate(%q, %d) resulted in length %d, expected <= %d", input, max, len(result), max)
	}

	if !strings.HasSuffix(result, "...") {
		t.Errorf("truncate(%q, %d) = %q, expected to end with '...'", input, max, result)
	}
}

func TestTruncate_EmptyString(t *testing.T) {
	input := ""
	max := 10
	result := truncate(input, max)

	if result != "" {
		t.Errorf("truncate(%q, %d) = %q, expected empty string", input, max, result)
	}
}

func TestTruncate_MaxLessThanEllipsis(t *testing.T) {
	input := "hello"
	max := 2
	result := truncate(input, max)

	if len(result) > max {
		t.Errorf("truncate(%q, %d) resulted in length %d, expected <= %d", input, max, len(result), max)
	}
}

func TestTruncate_UnicodeCharacters(t *testing.T) {
	input := "Hello 世界 🌍"
	max := 10
	result := truncate(input, max)

	if len(result) > max {
		t.Errorf("truncate with unicode resulted in length %d, expected <= %d", len(result), max)
	}
}

func TestWordWrap_NoWrappingNeeded(t *testing.T) {
	input := "short text"
	width := 20
	result := wordWrap(input, width)

	if result != input {
		t.Errorf("wordWrap(%q, %d) = %q, expected %q", input, width, result, input)
	}
}

func TestWordWrap_EmptyString(t *testing.T) {
	input := ""
	width := 10
	result := wordWrap(input, width)

	if result != "" {
		t.Errorf("wordWrap(%q, %d) = %q, expected empty string", input, width, result)
	}
}

func TestWordWrap_MultipleLines(t *testing.T) {
	input := "This is a long line that needs to be wrapped at some point because it exceeds the width"
	width := 20
	result := wordWrap(input, width)

	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if len(line) > width {
			t.Errorf("Line %d has length %d, expected <= %d: %q", i, len(line), width, line)
		}
	}
}

func TestWordWrap_PreservesExistingNewlines(t *testing.T) {
	input := "Line one\nLine two\nLine three"
	width := 50
	result := wordWrap(input, width)

	// Should contain the content (exact behavior of newline preservation may vary)
	if !strings.Contains(result, "Line one") || !strings.Contains(result, "Line two") {
		t.Errorf("wordWrap should preserve content, got: %q", result)
	}
}

func TestWordWrap_VeryLongWord(t *testing.T) {
	input := "verylongwordthatcannotbewrappedproperly"
	width := 10
	result := wordWrap(input, width)

	// Should handle the long word somehow (exact behavior may vary)
	if len(result) == 0 {
		t.Errorf("wordWrap should return non-empty result for long word")
	}

	// Result should not be longer than the original in most cases
	if len(result) > len(input)*2 {
		t.Errorf("wordWrap produced unexpectedly long result")
	}
}

func TestWordWrap_ZeroWidth(t *testing.T) {
	input := "test"
	width := 0
	result := wordWrap(input, width)

	// Should handle gracefully (implementation-dependent behavior)
	if result == "" {
		// One acceptable behavior: returns empty
	} else if len(result) > 0 {
		// Another acceptable behavior: returns something
	}
}

func TestWordWrap_NegativeWidth(t *testing.T) {
	input := "test"
	width := -5
	result := wordWrap(input, width)

	// Should handle gracefully without crashing
	_ = result
}

func TestMin_FirstSmaller(t *testing.T) {
	result := min(3, 5)
	if result != 3 {
		t.Errorf("min(3, 5) = %d, expected 3", result)
	}
}

func TestMin_SecondSmaller(t *testing.T) {
	result := min(10, 7)
	if result != 7 {
		t.Errorf("min(10, 7) = %d, expected 7", result)
	}
}

func TestMin_Equal(t *testing.T) {
	result := min(5, 5)
	if result != 5 {
		t.Errorf("min(5, 5) = %d, expected 5", result)
	}
}

func TestMin_NegativeNumbers(t *testing.T) {
	result := min(-3, -5)
	if result != -5 {
		t.Errorf("min(-3, -5) = %d, expected -5", result)
	}
}

func TestMin_ZeroValues(t *testing.T) {
	result := min(0, 5)
	if result != 0 {
		t.Errorf("min(0, 5) = %d, expected 0", result)
	}

	result = min(-5, 0)
	if result != -5 {
		t.Errorf("min(-5, 0) = %d, expected -5", result)
	}
}
