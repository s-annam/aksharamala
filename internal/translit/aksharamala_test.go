package translit

import (
	"testing"
)

func TestAksharamala_Transliterate(t *testing.T) {
	schemePath := "..\\..\\keymaps\\Hindi.aksj"
	aks, err := NewAksharamala(schemePath)
	if err != nil {
		t.Fatalf("Failed to initialize Aksharamala: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"namaste", "नमस्ते"},
		{"kk", "क्क"},
		{"ka", "क"},
		{"a1k", "अ१क"},
		{"", ""}, // Edge case: empty string
	}

	for _, test := range tests {
		output := aks.Transliterate(test.input)
		if output != test.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", test.input, test.expected, output)
		} else {
			t.Logf("For input '%s': output matches expected '%s'", test.input, test.expected)
		}
	}
}

// Test cases for parseVirama
func TestAksharamala_ParseVirama(t *testing.T) {
	tests := []struct {
		metadata     string
		expectedRune rune
		expectedMode string
		shouldError  bool
	}{
		{"0x094D,smart", '्', "smart", false},
		{"्,normal", '्', "normal", false},
		{"abcd,smart", 'a', "smart", false},
		{"0xZZZZ,smart", 0, "", true},
		{"no_comma", 0, "", true},
	}

	for _, test := range tests {
		r, mode, err := parseVirama(test.metadata)
		if test.shouldError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.metadata)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.metadata, err)
			continue
		}

		if r != test.expectedRune || mode != test.expectedMode {
			t.Errorf("For input '%s', expected rune '%c' and mode '%s', but got rune '%c' and mode '%s'",
				test.metadata, test.expectedRune, test.expectedMode, r, mode)
		}
	}
}
