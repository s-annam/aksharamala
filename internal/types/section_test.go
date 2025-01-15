package types

import (
	"testing"

	"aks.go/internal/core"
)

// TestAppendLHSToMapping tests the AppendLHSToMapping method of the Section struct.
// It verifies that new LHS values can be appended to existing mappings correctly.
func TestAppendLHSToMapping(t *testing.T) {
	section := NewSection()

	// Initial mapping to append to
	initialMapping := core.Mapping{
		LHS:     []string{"k"},
		RHS:     []string{"क"},
		Comment: "example mapping",
	}
	section.AddMapping(initialMapping.LHS, initialMapping.RHS, initialMapping.Comment)

	// Append a new LHS
	newLHS := "ka"
	section.AppendLHSToMapping(&initialMapping, newLHS)

	// Validate that the LHS was appended
	mappings := section.GetMappings()
	if len(mappings) == 0 {
		t.Fatal("Expected at least one mapping")
	}

	lastMapping := mappings[len(mappings)-1]
	if len(lastMapping.LHS) != 2 || lastMapping.LHS[1] != newLHS {
		t.Errorf("Expected LHS to include '%s', got %v", newLHS, lastMapping.LHS)
	}
}

// TestAddMappingAndGetMappings tests the retrieval of mappings from a Section.
// It verifies that mappings can be added and retrieved correctly.
func TestAddMappingAndGetMappings(t *testing.T) {
	section := NewSection()

	// Add mappings
	section.AddMapping([]string{"a"}, []string{"अ"}, "vowel")
	section.AddMapping([]string{"aa"}, []string{"आ"}, "long vowel")

	mappings := section.GetMappings()
	if len(mappings) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(mappings))
	}

	if mappings[0].LHS[0] != "a" || mappings[1].LHS[0] != "aa" {
		t.Errorf("Mappings do not match expected values: %v", mappings)
	}
}

// TestGetOrCreate tests the GetOrCreate method of the Section struct.
// It verifies that a section can be retrieved or created as needed.
func TestGetOrCreate(t *testing.T) {
	scheme := TransliterationScheme{
		Categories: make(map[string]Section),
	}

	categoryName := "vowels"
	section := GetOrCreate(scheme, categoryName)

	if section == nil {
		t.Fatalf("GetOrCreate returned nil")
	}

	// Ensure the section exists in the scheme
	_, exists := scheme.Categories[categoryName]
	if !exists {
		t.Errorf("Expected section '%s' to exist in scheme", categoryName)
	}
}

// TestValidateEmptyMappings tests the validation of empty mappings in a Section.
// It verifies that appropriate errors are raised for empty mappings.
func TestValidateEmptyMappings(t *testing.T) {
	section := NewSection()
	section.AddMapping([]string{}, []string{}, "empty mapping")
	err := section.Mappings.Validate("vowels", "test_scheme")
	if err == nil {
		t.Errorf("Expected error for empty mappings, got nil")
	}
}
