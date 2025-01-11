package types

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToCompactScheme(t *testing.T) {
	original := TransliterationScheme{
		Comments: []string{"Test comment"},
		Version:  "1.0",
		ID:       "test_id",
		Name:     "Test Scheme",
		License:  "Test License",
		Language: "Test Language",
		Scheme:   "Test Scheme",
		Metadata: Metadata{
			Virama:       "Test Virama",
			FontName:     "Test Font",
			FontSize:     12,
			IconEnabled:  "enabled.png",
			IconDisabled: "disabled.png",
		},
		Categories: map[string]Section{
			"test_category": {
				Comments: []string{"Category comment"},
				Mappings: []CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ"}},
				},
			},
		},
	}

	compact, err := ToCompactTransliterationScheme(original)
	if err != nil {
		t.Fatalf("Error converting to compact scheme: %v", err)
	}

	// Marshal the original and compact schemes to JSON
	originalJSON, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Error marshalling original scheme to JSON: %v", err)
	}

	compactJSON, err := json.Marshal(compact)
	if err != nil {
		t.Fatalf("Error marshalling compact scheme to JSON: %v", err)
	}

	// Unmarshal the JSON back to maps for comparison
	var originalMap map[string]interface{}
	var compactMap map[string]interface{}

	if err := json.Unmarshal(originalJSON, &originalMap); err != nil {
		t.Fatalf("Error unmarshalling original scheme JSON: %v", err)
	}

	if err := json.Unmarshal(compactJSON, &compactMap); err != nil {
		t.Fatalf("Error unmarshalling compact scheme JSON: %v", err)
	}

	// Remove the categories field for comparison
	delete(originalMap, "categories")
	delete(compactMap, "categories")

	// Compare the remaining fields
	if !equalMaps(originalMap, compactMap) {
		t.Fatalf("Original and compact schemes do not match")
	}
}

// Helper function to compare two maps
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !equalValues(v, bv) {
			return false
		}
	}
	return true
}

// Helper function to compare two values
func equalValues(a, b interface{}) bool {
	switch a := a.(type) {
	case map[string]interface{}:
		b, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		return equalMaps(a, b)
	case []interface{}:
		b, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if !equalValues(a[i], b[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

func TestMappingFormat(t *testing.T) {
	// Create a test scheme with multiple mappings
	scheme := TransliterationScheme{
		ID:       "test",
		Name:     "Test Scheme",
		Language: "Test",
		Scheme:   "Test",
		Categories: map[string]Section{
			"test_category": {
				Mappings: []CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"A"}},
					{LHS: []string{"b"}, RHS: []string{"B"}},
					{LHS: []string{"c", "C"}, RHS: []string{"क"}},
				},
			},
		},
	}

	// Convert to compact scheme
	compact, err := ToCompactTransliterationScheme(scheme)
	if err != nil {
		t.Fatalf("Failed to convert to compact scheme: %v", err)
	}

	// Get the raw JSON for the test category
	categoryJSON, ok := compact.Categories["test_category"]
	if !ok {
		t.Fatal("Test category not found in compact scheme")
	}

	// Convert to string for inspection
	jsonStr := string(categoryJSON)
	t.Logf("JSON output:\n%s", jsonStr) // Add this for debugging

	// Test formatting expectations
	t.Run("Array brackets placement", func(t *testing.T) {
		if !strings.HasPrefix(jsonStr, "[") {
			t.Error("Expected array to start with '[', got:", jsonStr[:1])
		}
		if !strings.HasSuffix(jsonStr, "]") {
			t.Error("Expected array to end with ']'")
		}
	})

	t.Run("One mapping per line", func(t *testing.T) {
		lines := strings.Split(jsonStr, "\n")
		expectedLines := len(scheme.Categories["test_category"].Mappings) + 2 // +2 for opening and closing brackets
		if len(lines) != expectedLines {
			t.Errorf("Expected %d lines, got %d", expectedLines, len(lines))
		}
	})

	t.Run("Mapping format", func(t *testing.T) {
		lines := strings.Split(jsonStr, "\n")
		for i, line := range lines {
			// Skip first and last lines (brackets)
			if i > 0 && i < len(lines)-1 {
				if !strings.HasPrefix(line, "      {") {
					t.Errorf("Line %d should start with '      {', got: %s", i, line)
				}
				if i < len(lines)-2 && !strings.HasSuffix(line, "},") {
					t.Errorf("Line %d should end with '},', got: %s", i, line)
				}
				if i == len(lines)-2 && !strings.HasSuffix(line, "}") {
					t.Errorf("Last mapping line should end with '}', got: %s", line)
				}
			}
		}
	})
}

// Test entire JSON output
func TestFullJSONOutput(t *testing.T) {
	scheme := TransliterationScheme{
		ID:       "test",
		Name:     "Test Scheme",
		Language: "Test",
		Scheme:   "Test",
		Categories: map[string]Section{
			"consonants": {
				Mappings: []CategoryEntry{
					{LHS: []string{"k"}, RHS: []string{"क"}},
					{LHS: []string{"kh"}, RHS: []string{"ख"}},
				},
			},
		},
	}

	compact, err := ToCompactTransliterationScheme(scheme)
	if err != nil {
		t.Fatalf("Failed to convert to compact scheme: %v", err)
	}

	// Encode with indentation
	var buf strings.Builder
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(compact); err != nil {
		t.Fatalf("Failed to encode compact scheme: %v", err)
	}

	output := buf.String()
	t.Logf("Full JSON output:\n%s", output)

	// Test the structure
	if !strings.Contains(output, `"consonants": [`) {
		t.Error("Missing expected category opening")
	}

	// Split into lines and check each line
	lines := strings.Split(output, "\n")
	foundMapping := false
	for i := 0; i < len(lines)-1; i++ {
		// If we find a [ character, the next line should be a mapping
		if strings.Contains(lines[i], `": [`) {
			nextLine := lines[i+1]
			leadingSpaces := len(nextLine) - len(strings.TrimLeft(nextLine, " "))
			if leadingSpaces != 6 {
				t.Errorf("Expected 6 leading spaces for mapping, got %d: %s",
					leadingSpaces, nextLine)
			}
			foundMapping = true
			break
		}
	}

	if !foundMapping {
		t.Error("No category mappings found in output")
	}
}
