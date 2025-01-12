package translit

import (
	"testing"

	"aks.go/internal/types"
)

// TestViramaBasic focuses on core virama behavior between consonants
// without handling special cases like .r
func TestViramaBasic(t *testing.T) {
	// Define test scheme with virama handling
	scheme := &types.TransliterationScheme{
		Metadata: types.Metadata{
			// Use the same virama character as in Hindi.aksj
			Virama: "0x094D, smart", // DEVANAGARI SIGN VIRAMA with smart mode
		},
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"k"}, RHS: []string{"क"}},
					{LHS: []string{"r"}, RHS: []string{"र"}},
					{LHS: []string{"m"}, RHS: []string{"म"}},
					{LHS: []string{"p"}, RHS: []string{"प"}},
					{LHS: []string{"n"}, RHS: []string{"न"}},
					{LHS: []string{"s"}, RHS: []string{"स"}},
					{LHS: []string{"t"}, RHS: []string{"त"}},
				},
			},
			"vowels": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ", ""}},    // Base vowel with empty matra
					{LHS: []string{"aa"}, RHS: []string{"आ", "ा"}},  // Long A with matra
					{LHS: []string{"e"}, RHS: []string{"ए", "े"}},   // E with matra
					{LHS: []string{"i"}, RHS: []string{"इ", "ि"}},   // I with matra
					{LHS: []string{"R^i"}, RHS: []string{"ऋ", "ृ"}}, // Vocalic R with matra
				},
			},
		},
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "basic virama between consonants (namaste)",
			input: "namaste",
			want:  "नमस्ते",
		},
		{
			name:  "virama with final consonant (karma)",
			input: "karma",
			want:  "कर्म",
		},
		{
			name:  "consonant + r + vowel combination (kripa)",
			input: "kripa",
			want:  "कृप",
		},
		{
			name:  "multiple viramas (matsya)",
			input: "matsya",
			want:  "मत्स्य",
		},
		{
			name:  "r as second consonant (prana)",
			input: "prana",
			want:  "प्राण",
		},
		{
			name:  "standalone consonant with implicit a (ka)",
			input: "ka",
			want:  "क",
		},
		{
			name:  "consonant cluster with different vowel (kri)",
			input: "kri",
			want:  "कृ",
		},
		{
			name:  "three consonant cluster (stri)",
			input: "stri",
			want:  "स्त्री",
		},
	}

	aks, err := NewAksharamala(scheme)
	if err != nil {
		t.Fatalf("Failed to create Aksharamala: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aks.Transliterate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transliterate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Transliterate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestViramaMode tests different virama modes from the C++ implementation
func TestViramaMode(t *testing.T) {
	testModes := []struct {
		name       string
		viramaStr  string // Virama configuration string
		testInputs []struct {
			input string
			want  string
		}
	}{
		{
			name:      "smart mode",
			viramaStr: "0x094D, smart",
			testInputs: []struct {
				input string
				want  string
			}{
				{"karma", "कर्म"},
				{"kripa", "कृप"},
			},
		},
		{
			name:      "normal mode",
			viramaStr: "0x094D, normal",
			testInputs: []struct {
				input string
				want  string
			}{
				{"karma", "कर्म"},
				{"kripa", "क्रिप"},
			},
		},
		{
			name:      "double mode",
			viramaStr: "0x094D, double",
			testInputs: []struct {
				input string
				want  string
			}{
				{"kk", "क्क"},    // Double consonant
				{"kkaa", "क्का"}, // Double consonant with vowel
			},
		},
		{
			name:      "repeat mode",
			viramaStr: "0x094D, repeat",
			testInputs: []struct {
				input string
				want  string
			}{
				{"kk", "क्क"},    // Repeated consonant
				{"rr", "र्र"},    // Different repeated consonant
				{"ka", "क"},      // Non-repeated should work normally
				{"krr", "क्र्र"}, // Multiple repeats
			},
		},
	}

	baseScheme := &types.TransliterationScheme{
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"k"}, RHS: []string{"क"}},
					{LHS: []string{"r"}, RHS: []string{"र"}},
				},
			},
			"vowels": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ", ""}},
					{LHS: []string{"aa"}, RHS: []string{"आ", "ा"}},
					{LHS: []string{"i"}, RHS: []string{"इ", "ि"}},
				},
			},
		},
	}

	for _, mode := range testModes {
		t.Run(mode.name, func(t *testing.T) {
			// Create a copy of base scheme with specific virama mode
			scheme := *baseScheme
			scheme.Metadata.Virama = mode.viramaStr

			aks, err := NewAksharamala(&scheme)
			if err != nil {
				t.Fatalf("Failed to create Aksharamala: %v", err)
			}

			for _, tc := range mode.testInputs {
				got, err := aks.Transliterate(tc.input)
				if err != nil {
					t.Errorf("Transliterate(%q) error = %v", tc.input, err)
					continue
				}
				if got != tc.want {
					t.Errorf("Transliterate(%q) = %q, want %q", tc.input, got, tc.want)
				}
			}
		})
	}
}
