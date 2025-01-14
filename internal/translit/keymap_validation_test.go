package translit

import (
	"encoding/json"
	"testing"

	"aks.go/internal/core"
	"aks.go/internal/types"
)

// TestValidateKeymap tests the validation of a keymap.
// It verifies that a valid keymap is accepted and invalid keymaps are rejected.
func TestValidateKeymap(t *testing.T) {
	validKeymap := types.TransliterationScheme{
		ID:       "hindi",
		Name:     "Hindi Keymap",
		Language: "Hindi",
		Scheme:   "ITRANS",
		Categories: map[string]types.Section{
			"vowels": {
				Mappings: core.NewMappings([]core.Mapping{
					{LHS: []string{"a"}, RHS: []string{"अ"}},
				}),
			},
		},
		Metadata: types.Metadata{
			Virama: "0x094D",
		},
	}

	tests := []struct {
		name         string
		keymap       types.TransliterationScheme
		expectsError bool
	}{
		{
			"Valid keymap",
			validKeymap,
			false,
		},
		{
			"Missing ID",
			types.TransliterationScheme{
				Categories: validKeymap.Categories,
				Metadata:   validKeymap.Metadata,
			},
			true,
		},
		{
			"Empty Categories",
			types.TransliterationScheme{
				ID:         "hindi",
				Categories: map[string]types.Section{},
				Metadata:   validKeymap.Metadata,
			},
			true,
		},
		{
			"Invalid LHS",
			types.TransliterationScheme{
				ID: "hindi",
				Categories: map[string]types.Section{
					"vowels": {
						Mappings: core.NewMappings(
							[]core.Mapping{
								{LHS: []string{}, RHS: []string{"अ"}},
							}),
					},
				},
				Metadata: validKeymap.Metadata,
			},
			true,
		},
		{
			"Invalid RHS",
			types.TransliterationScheme{
				ID: "hindi",
				Categories: map[string]types.Section{
					"vowels": {
						Mappings: core.NewMappings([]core.Mapping{
							{LHS: []string{"a"}, RHS: []string{}},
						}),
					},
				},
				Metadata: validKeymap.Metadata,
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.keymap.Validate()
			if test.expectsError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !test.expectsError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
		})
	}
}

// TestInvalidKeymapJSON tests the handling of invalid JSON for keymaps.
// It verifies that the system correctly identifies and rejects invalid keymap JSON.
func TestInvalidKeymapJSON(t *testing.T) {
	invalidJSON := `{
		"id": "hindi",
		"categories": {
			"vowels": {
				"mappings": [
					{"lhs": ["a"]}
				]
			}
		},
		"metadata": {}
	}`

	var keymap types.TransliterationScheme
	err := json.Unmarshal([]byte(invalidJSON), &keymap)
	if err == nil {
		t.Errorf("Expected error for invalid JSON but got none")
	} else {
		t.Logf("Correctly caught JSON error: %v", err)
	}
}
