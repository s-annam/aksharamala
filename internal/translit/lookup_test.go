package translit

import (
	"testing"
)

func TestBuildLookupTable(t *testing.T) {
	scheme := &TransliterationScheme{
		Categories: map[string][]Mapping{
			"consonants": {
				{LHS: []string{"k"}, RHS: []string{"क"}},
				{LHS: []string{"kh"}, RHS: []string{"ख"}},
			},
			"vowels": {
				{LHS: []string{"a"}, RHS: []string{"अ"}},
			},
		},
	}

	table := BuildLookupTable(scheme)

	// Test valid lookups
	tests := []struct {
		input    string
		expected LookupResult
	}{
		{"k", LookupResult{Output: "क", Category: "consonants"}},
		{"a", LookupResult{Output: "अ", Category: "vowels"}},
	}

	for _, test := range tests {
		result := table.Lookup(test.input)
		if result != test.expected {
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
