package translit

import (
	"fmt"
	"testing"

	"aks.go/internal/keymap"
	"aks.go/internal/types"
)

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

// Test cases for parseVirama
func TestAksharamala_ParseVirama(t *testing.T) {
	tests := []struct {
		metadata       string
		expectedVirama string
		expectedMode   types.ViramaMode
		shouldError    bool
	}{
		{"0x094D,smart", "्", types.SmartMode, false},
		{"्,normal", "्", types.NormalMode, false},
		{"abcd,smart", "abcd", types.SmartMode, false},
		{"0xZZZZ,smart", "", types.UnknownMode, true},
		{"no_comma", "", types.UnknownMode, true},
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

		if r != test.expectedVirama || mode != test.expectedMode {
			t.Errorf("For input '%s', expected '%s' and mode '%s', but got '%s' and mode '%s'",
				test.metadata, test.expectedVirama, test.expectedMode, r, mode)
		}
	}
}
