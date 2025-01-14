package translit

import (
	"fmt"
	"testing"

	"aks.go/internal/keymap"
)

// TestAksharamala_Transliterate tests the Transliterate method of the Aksharamala struct.
// It verifies that the transliteration works correctly for various input cases.
// The test loads keymaps and checks the output against expected results.
func TestAksharamala_Transliterate(t *testing.T) {
	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps("../../keymaps"); err != nil {
		fmt.Printf("Failed to load keymaps: %v\n", err)
		return
	}

	aks := NewAksharamala(store)

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

	aks.SetActiveKeymap("hindi")
	for _, test := range tests {
		output := aks.Transliterate(test.input)
		if output != test.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", test.input, test.expected, output)
		} else {
			t.Logf("For input '%s': output matches expected '%s'", test.input, test.expected)
		}
	}
}
