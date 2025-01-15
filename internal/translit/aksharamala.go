package translit

import (
	"fmt"
	"strings"

	"aks.go/internal/core"
	"aks.go/internal/keymap"
	"aks.go/internal/types"
)

// Aksharamala represents a transliteration engine that uses a keymap store
// to perform transliteration operations. It maintains the active transliteration
// scheme and manages the context for transliteration processes.
type Aksharamala struct {
	keymapStore  *keymap.KeymapStore
	activeScheme *types.TransliterationScheme
	context      *types.Context
	virama       string
	viramaMode   types.ViramaMode
}

// NewAksharamala initializes a new Aksharamala instance.
func NewAksharamala(store *keymap.KeymapStore) *Aksharamala {
	return &Aksharamala{
		keymapStore: store,
		context:     types.NewContext(),
	}
}

// SetActiveKeymap sets the active keymap by ID for transliteration.
// Returns an error if the keymap ID is not found.
func (a *Aksharamala) SetActiveKeymap(id string) error {
	scheme, exists := a.keymapStore.GetKeymap(id)
	if !exists {
		return fmt.Errorf("keymap with ID '%s' not found", id)
	}
	a.activeScheme = &scheme
	a.virama, a.viramaMode, _ = types.ParseVirama(scheme.Metadata.Virama) // Assuming validation is done earlier
	return nil
}

// TransliterateWithKeymap performs transliteration for consonants, vowels, and mixed input
// using the active keymap. Returns the transliterated string and an error if any occurs.
func (a *Aksharamala) TransliterateWithKeymap(id, input string) (string, error) {
	if err := a.SetActiveKeymap(id); err != nil {
		return "", err
	}
	return a.Transliterate(input), nil
}

// Transliterate performs transliteration for the input string using the active scheme.
// Returns the transliterated string.
func (a *Aksharamala) Transliterate(input string) string {
	// Reset context for a clean state
	a.context = types.NewContext()

	var result strings.Builder
	length := len(input)
	for i := 0; i < length; {
		foundMatch := false

		// Greedily match substrings to find the longest match
		for j := length - i; j > 0; j-- {
			if i+j <= length {
				combination := input[i : i+j]
				lookupResult := a.lookup(combination)

				if lookupResult.Output != "" {
					if lookupResult.Output == "\x00" && a.context.LatestLookup.Category == "consonants" {
						a.context.LatestLookup = lookupResult
						i += j // Move the index forward by the length of the match
						foundMatch = true
						break
					}

					// Apply virama if needed
					if a.shouldApplyVirama(lookupResult.Output) {
						result.WriteString(a.virama)
					}

					result.WriteString(lookupResult.Output)
					a.context.LatestLookup = lookupResult
					i += j // Move the index forward by the length of the match
					foundMatch = true
					break // Exit the loop once a match is found
				}
			}
		}

		// If no match was found, process the current character
		if !foundMatch {
			char := string(input[i])
			lookupResult := a.lookup(char)
			if lookupResult.Output != "" {
				result.WriteString(lookupResult.Output)
				a.context.LatestLookup = lookupResult
			} else {
				result.WriteString(char) // Handle unmatched characters
				a.context.LatestLookup = core.LookupResult{Output: char, Category: "other"}
			}
			i++ // Move to the next character
		}
	}
	return result.String()
}

// shouldApplyVirama determines if a virama should be inserted before the current character.
// Returns true if a virama should be applied, false otherwise.
func (a *Aksharamala) shouldApplyVirama(nextOutput string) bool {
	if a.context.LatestLookup.Category != "consonants" || a.getCategory(nextOutput) != "consonants" {
		return false
	}

	switch a.viramaMode {
	case types.SmartMode:
		return true
	case types.NormalMode:
		return true
	case types.UnknownMode:
		return a.context.LatestLookup.Output == nextOutput
	}

	return false
}

// lookup finds the transliteration for the given string.
// Returns the LookupResult for the character.
func (a *Aksharamala) lookup(combination string) core.LookupResult {
	for category, section := range a.activeScheme.Categories {
		for _, mapping := range section.Mappings.Entries() {
			for _, lhs := range mapping.LHS {
				if lhs == combination {
					// Use matra (RHS[1]) if the previous character is a consonant
					if category == "vowels" && a.context.LatestLookup.Category == "consonants" && len(mapping.RHS) > 1 {
						return core.LookupResult{Output: mapping.RHS[1], Category: category}
					}
					return core.LookupResult{Output: mapping.RHS[0], Category: category} // Use full form otherwise
				}
			}
		}
	}
	return core.LookupResult{Output: "", Category: "other"} // No match found
}

// getCategory determines the category of the output character.
// Returns the category as a string.
func (a *Aksharamala) getCategory(output string) string {
	for category, section := range a.activeScheme.Categories {
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
