package types

import (
	"testing"

	"aks.go/internal/core"
)

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
	err := section.AppendLHSToMapping(&initialMapping, newLHS)
	if err != nil {
		t.Fatalf("AppendLHSToMapping failed: %v", err)
	}

	// Validate that the LHS was appended
	lastMapping := section.GetMappings()[len(section.GetMappings())-1]
	if len(lastMapping.LHS) != 2 || lastMapping.LHS[1] != newLHS {
		t.Errorf("Expected LHS to include '%s', got %v", newLHS, lastMapping.LHS)
	}
}

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

func TestValidateEmptyMappings(t *testing.T) {
	section := NewSection()
	section.AddMapping([]string{}, []string{}, "empty mapping")
	err := section.Mappings.Validate("vowels", "test_scheme")
	if err == nil {
		t.Errorf("Expected error for empty mappings, got nil")
	}
}
