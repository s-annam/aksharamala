package translit

import (
	"reflect"
	"testing"

	"aks.go/internal/core"
	"aks.go/internal/types"
)

// TestBuildLookupTable tests the building of a lookup table from a TransliterationScheme.
// It verifies that the lookup table is constructed correctly based on the scheme.
func TestBuildLookupTable(t *testing.T) {
	scheme := &types.TransliterationScheme{
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: core.NewMappings([]core.Mapping{
					{LHS: []string{"k"}, RHS: []string{"क"}},
					{LHS: []string{"kh"}, RHS: []string{"ख"}},
				}),
			},
			"vowels": {
				Mappings: core.NewMappings([]core.Mapping{
					{LHS: []string{"a"}, RHS: []string{"अ"}},
				}),
			},
		},
	}

	table := BuildLookupTable(scheme)

	// Test valid lookups
	tests := []struct {
		input    string
		expected core.LookupResult
	}{
		{"k", core.LookupResult{Output: "क", Category: "consonants", Found: true}},
		{"a", core.LookupResult{Output: "अ", Category: "vowels", Found: true}},
	}

	for _, test := range tests {
		result := table.Lookup(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input '%s': expected %+v, got %+v", test.input, test.expected, result)
		} else {
			t.Logf("For input '%s': got expected result %+v", test.input, result)
		}
	}

	// Test invalid lookup
	result := table.Lookup("z")
	if result.Output != "" || result.Category != "other" {
		t.Errorf("For input 'z': expected empty output and 'other' category, got %+v", result)
	}
}
