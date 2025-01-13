// Final fixes for vowel handling in Stage 4
package translit

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"aks.go/internal/types"
)

type Aksharamala struct {
	scheme     *types.TransliterationScheme
	context    *Context
	virama     string
	viramaMode string
}

func splitAndTrim(s string) []string {
	// Split the string by commas and trim white spaces
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ','
	})

	// Trim white spaces from each part
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}

func parseVirama(metadata string) (string, string, error) {
	// Split virama metadata entry into two parts
	parts := splitAndTrim(metadata)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid or missing virama metadata: %s", metadata)
	}
	viramaString := parts[0]
	viramaMode := parts[1]

	// Handle hexadecimal Unicode value, if present
	if strings.HasPrefix(viramaString, "0x") {
		codePoint, err := strconv.ParseUint(viramaString[2:], 16, 32)
		if err != nil {
			return "", "", fmt.Errorf("invalid Unicode code point: %v", err)
		}
		return string(rune(codePoint)), viramaMode, nil
	} else if viramaString == "" {
		return "", "", fmt.Errorf("empty virama: %s", metadata)
	}

	return viramaString, viramaMode, nil
}

// NewAksharamala initializes a new Aksharamala instance.
func NewAksharamala(schemePath string) (*Aksharamala, error) {
	data, err := os.ReadFile(schemePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scheme: %w", err)
	}

	scheme := &types.TransliterationScheme{}
	if err := json.Unmarshal(data, scheme); err != nil {
		return nil, fmt.Errorf("failed to parse scheme: %w", err)
	}

	viramaRune, viramaMode, err := parseVirama(scheme.Metadata.Virama)
	if err != nil {
		return nil, fmt.Errorf("invalid virama: %w", err)
	}

	return &Aksharamala{
		scheme:     scheme,
		context:    NewContext(),
		virama:     viramaRune,
		viramaMode: viramaMode,
	}, nil
}

// Transliterate performs transliteration for consonants, vowels, and mixed input.
func (a *Aksharamala) Transliterate(input string) string {
	// Reset context for a clean state
	a.context = NewContext()

	var result strings.Builder
	for _, char := range input {
		lookupResult := a.lookup(string(char))
		if lookupResult.Output == "\x00" && a.context.LatestLookup.Category == "consonants" {
			a.context.LatestLookup = lookupResult
			continue
		}

		if lookupResult.Output != "" {
			// Handle virama based on mode and context
			if a.shouldApplyVirama(lookupResult.Output) {
				result.WriteString(a.virama)
			}

			// Update context and write output
			a.context.LatestLookup = lookupResult
			result.WriteString(lookupResult.Output)
		} else {
			// For unmatched characters, treat as "other"
			result.WriteString(string(char))
			a.context.LatestLookup = LookupResult{Output: string(char), Category: "other"}
		}
	}
	return result.String()
}

// shouldApplyVirama determines if a virama should be inserted before the current character.
func (a *Aksharamala) shouldApplyVirama(nextOutput string) bool {
	if a.context.LatestLookup.Category != "consonants" || a.getCategory(nextOutput) != "consonants" {
		return false
	}

	switch a.viramaMode {
	case "smart":
		return true
	case "normal":
		return true
	case "double":
		return a.context.LatestLookup.Output == nextOutput
	case "repeat":
		return a.context.LatestLookup.Output == nextOutput
	}

	return false
}

// lookup finds the transliteration for a single character.
func (a *Aksharamala) lookup(char string) LookupResult {
	for category, section := range a.scheme.Categories {
		for _, mapping := range section.Mappings.Entries() {
			for _, lhs := range mapping.LHS {
				if lhs == char {
					// Use matra (RHS[1]) if the previous character is a consonant
					if category == "vowels" && a.context.LatestLookup.Category == "consonants" && len(mapping.RHS) > 1 {
						return LookupResult{Output: mapping.RHS[1], Category: category}
					}
					return LookupResult{Output: mapping.RHS[0], Category: category} // Use full form otherwise
				}
			}
		}
	}
	return LookupResult{Output: "", Category: "other"} // No match found
}

// getCategory determines the category of the output character.
func (a *Aksharamala) getCategory(output string) string {
	for category, section := range a.scheme.Categories {
		for _, mapping := range section.Mappings.Entries() {
			for _, rhs := range mapping.RHS {
				if rhs == output {
					return category
				}
			}
		}
	}
	return "other"
}
