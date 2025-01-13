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
