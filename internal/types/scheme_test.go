package types

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToCompactScheme(t *testing.T) {
	original := TransliterationScheme{
		Comments: []string{"Test comment"},
		Version:  "2025.1",
		ID:       "test_id",
		Name:     "Test Scheme",
		License:  "AGPL",
		Language: "Test Language",
		Scheme:   "Unicode",
		Metadata: Metadata{
			Virama:       "Test Virama",
			FontName:     "Test Font",
			FontSize:     12,
			IconEnabled:  "enabled.png",
			IconDisabled: "disabled.png",
		},
		Categories: map[string]Section{
			"vowels": {
				Comments: []string{"Category comment"},
				Mappings: []CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ"}, Comment: "=*= Vowel A =*="},
					{LHS: []string{"aa"}, RHS: []string{"आ"}, Comment: " Vowel AA "},
				},
			},
		},
	}

	compact, err := ToCompactTransliterationScheme(original)
	if err != nil {
		t.Fatalf("Error converting to compact scheme: %v", err)
	}

	// Test version and license
	if compact.Version != "2025.1" {
		t.Errorf("Expected version '2025.1', got '%s'", compact.Version)
	}
	if compact.License != "AGPL" {
		t.Errorf("Expected license 'AGPL', got '%s'", compact.License)
	}

	// Test categories
	if len(compact.Categories) != len(original.Categories) {
		t.Fatalf("Mismatch in number of categories: expected %d, got %d",
			len(original.Categories), len(compact.Categories))
	}

	// Check category serialization
	categoryJSON, ok := compact.Categories["vowels"]
	if !ok {
		t.Fatal("Vowels category not found in compact scheme")
	}

	var mappings []CategoryEntry
	if err := json.Unmarshal(categoryJSON, &mappings); err != nil {
		t.Fatalf("Error unmarshalling category JSON: %v", err)
	}

	// Test normalized comments
	expectedComments := []string{
		"Vowel A",
		"Vowel AA",
	}
	for i, mapping := range mappings {
		if mapping.Comment != expectedComments[i] {
			t.Errorf("Expected comment '%s', got '%s'",
				expectedComments[i], mapping.Comment)
		}
	}
}

func TestFullJSONOutput(t *testing.T) {
	scheme := TransliterationScheme{
		Comments: []string{"Generated from Hindi.akt. Distributed under AGPL."},
		Version:  "2025.1",
		ID:       "hindi",
		Name:     "Hindi Transliteration",
		License:  "AGPL",
		Language: "Hindi",
		Scheme:   "Unicode",
		Categories: map[string]Section{
			"consonants": {
				Mappings: []CategoryEntry{
					{LHS: []string{"k"}, RHS: []string{"क"}, Comment: "Consonant K"},
					{LHS: []string{"kh"}, RHS: []string{"ख"}, Comment: "Consonant KH"},
				},
			},
		},
	}

	compact, err := ToCompactTransliterationScheme(scheme)
	if err != nil {
		t.Fatalf("Failed to convert to compact scheme: %v", err)
	}

	// Serialize to JSON
	var buf strings.Builder
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(compact); err != nil {
		t.Fatalf("Failed to encode compact scheme: %v", err)
	}

	output := buf.String()
	t.Logf("Full JSON output:\n%s", output)

	// Check basic structure
	if !strings.Contains(output, `"version": "2025.1"`) {
		t.Error("Version field missing or incorrect")
	}
	if !strings.Contains(output, `"license": "AGPL"`) {
		t.Error("License field missing or incorrect")
	}
	if !strings.Contains(output, `"consonants": [`) {
		t.Error("Consonants category missing or incorrectly formatted")
	}
}
