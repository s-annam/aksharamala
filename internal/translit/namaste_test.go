package translit

import (
	"testing"

	"aks.go/internal/types"
)

func TestDebugNamaste(t *testing.T) {
	// Define test scheme with necessary mappings
	scheme := &types.TransliterationScheme{
		Metadata: types.Metadata{
			Virama: "0x094D, smart", // DEVANAGARI SIGN VIRAMA
		},
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"n"}, RHS: []string{"न"}},
					{LHS: []string{"m"}, RHS: []string{"म"}},
					{LHS: []string{"s"}, RHS: []string{"स"}},
					{LHS: []string{"t"}, RHS: []string{"त"}},
				},
			},
			"vowels": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ", ""}},
					{LHS: []string{"e"}, RHS: []string{"ए", "े"}},
				},
			},
		},
	}

	aks, err := NewAksharamala(scheme)
	if err != nil {
		t.Fatalf("Failed to create Aksharamala: %v", err)
	}

	word := "namaste"
	t.Logf("Testing full word: '%s'", word)

	// First, try transliterating the full word
	result, err := aks.Transliterate(word)
	if err != nil {
		t.Fatalf("Failed to transliterate word: %v", err)
	}
	t.Logf("Full word result: '%s', expected: 'नमस्ते'", result)

	if result != "नमस्ते" {
		t.Errorf("Full word transliteration failed. Got '%s', want 'नमस्ते'", result)
	}

	// Now test each character to debug the process
	t.Log("\nTesting character by character:")
	var output string
	for i, r := range word {
		// Show buffer state before processing
		t.Logf("Before char '%c': Buffer contains: '%s'", r, aks.buffer.String())

		result, err := aks.Transliterate(string(r))
		if err != nil {
			t.Errorf("Failed to transliterate '%c': %v", r, err)
			continue
		}

		output += result
		t.Logf("After char '%c': Got '%s', Buffer: '%s', LastChar: %c",
			r, result, aks.buffer.String(), aks.context.LastChar)

		// Special logging for 's'
		if r == 's' {
			nextChar := rune(0)
			if i+1 < len(word) {
				nextChar = rune(word[i+1])
			}
			t.Logf("Processing 's': Next char is '%c', Buffer after: '%s'",
				nextChar, aks.buffer.String())
		}
	}

	t.Logf("Character by character result: '%s'", output)
}
