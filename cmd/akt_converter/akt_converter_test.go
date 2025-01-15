package main

import (
	"testing"

	"aks.go/internal/core"
	"aks.go/internal/types"
	"aks.go/logger"
)

func init() {
	logger.InitLogger(false)
}

// TestConvertToCompactSchemeSectionMismatch verifies that when updating an existing scheme,
// if a mapping is found by LHS but in a different section, it updates the existing mapping
// in its original section rather than creating a new one.
func TestConvertToCompactSchemeSectionMismatch(t *testing.T) {
	// Create existing scheme with a mapping in "consonants" section
	existingScheme := types.TransliterationScheme{
		Version:  "2025.1",
		ID:       "test",
		Name:     "Test Scheme",
		Language: "Test",
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: core.NewMappings([]core.Mapping{
					{
						LHS:     []string{"x", "ksh"},
						RHS:     []string{"क्ष"},
						Comment: "original comment",
					},
				}),
			},
		},
	}

	// Create input scheme with same mapping but in "other_consonants" section
	inputScheme := types.TransliterationScheme{
		Version:  "2025.1",
		ID:       "test",
		Name:     "Test Scheme",
		Language: "Test",
		Categories: map[string]types.Section{
			"other_consonants": {
				Mappings: core.NewMappings([]core.Mapping{
					{
						LHS:     []string{"ksh"},
						RHS:     []string{"क्ष"},
						Comment: "updated comment",
					},
				}),
			},
		},
	}

	// Convert with update mode
	result, err := convertToCompactScheme(inputScheme, "test.akt", &existingScheme)
	if err != nil {
		t.Fatalf("convertToCompactScheme failed: %v", err)
	}

	// Convert back to regular scheme
	finalScheme, err := types.FromCompactTransliterationScheme(result)
	if err != nil {
		t.Fatalf("Failed to convert from compact scheme: %v", err)
	}

	// Verify that there is still only one "consonants" section
	consonantsSection, exists := finalScheme.Categories["consonants"]
	if !exists {
		t.Error("Original 'consonants' section not found in result")
	}

	// Verify that "other_consonants" section was not created
	if _, exists := finalScheme.Categories["other_consonants"]; exists {
		t.Error("Unexpected 'other_consonants' section was created")
	}

	// Verify that the mapping was updated in the original section
	mappings := consonantsSection.Mappings.All()
	if len(mappings) != 1 {
		t.Errorf("Expected 1 mapping in consonants section, got %d", len(mappings))
	}

	mapping := mappings[0]
	// Verify the mapping was updated but kept its original LHS array
	if len(mapping.LHS) != 2 || mapping.LHS[0] != "x" || mapping.LHS[1] != "ksh" {
		t.Errorf("Expected LHS to be ['x', 'ksh'], got %v", mapping.LHS)
	}
	if len(mapping.RHS) != 1 || mapping.RHS[0] != "क्ष" {
		t.Errorf("Expected RHS to be ['क्ष'], got %v", mapping.RHS)
	}
	if mapping.Comment != "updated comment" {
		t.Errorf("Expected comment to be 'updated comment', got %s", mapping.Comment)
	}
}
