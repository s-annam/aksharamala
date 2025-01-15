package core

import "testing"

// TestMappingsAddAndLookup tests the Add and Lookup methods of the Mappings struct.
// It verifies that mappings can be added and subsequently looked up correctly.
// The test checks both an existing mapping and a non-existing mapping.
func TestMappingsAddAndLookup(t *testing.T) {
	mappings := &Mappings{}
	mappings.Add([]string{"a"}, []string{"अ"}, "Vowel A")

	// Test lookup for an existing mapping
	output, exists := mappings.Lookup("a")
	if !exists || output != "अ" {
		t.Errorf("Expected 'अ', got '%s'", output)
	}

	// Test lookup for a non-existing mapping
	_, exists = mappings.Lookup("z")
	if exists {
		t.Errorf("Expected no match for 'z', but found one")
	}
}

func TestMappingsValidate(t *testing.T) {
	mappings := &Mappings{}
	mappings.Add([]string{"a"}, []string{"अ"}, "Vowel A")
	mappings.Add([]string{"i"}, []string{"इ"}, "Vowel I")

	// Test valid mappings
	err := mappings.Validate("test_category", "test_schemeID")
	if err != nil {
		t.Errorf("Validation failed for valid mappings: %v", err)
	}

	// Test invalid mapping (empty LHS)
	mappings.Add([]string{}, []string{"उ"}, "Invalid Mapping")
	err = mappings.Validate("test_category", "test_schemeID")
	if err == nil {
		t.Errorf("Expected validation error for empty LHS, but got none")
	}
}

func TestMappingsAll(t *testing.T) {
	mappings := &Mappings{}
	mappings.Add([]string{"a"}, []string{"अ"}, "Vowel A")
	mappings.Add([]string{"i"}, []string{"इ"}, "Vowel I")

	allMappings := mappings.All()
	if len(allMappings) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(allMappings))
	}
}
