package translit

import (
	"fmt"
	"strings"

	"aks.go/internal/types"
)

// Add viramaHandler to Aksharamala struct
type Aksharamala struct {
	scheme        *types.TransliterationScheme
	buffer        *Buffer
	context       *Context
	categoryMaps  map[string]map[string][]string
	viramaHandler *ViramaHandler
}

// Update NewAksharamala to initialize viramaHandler
func NewAksharamala(scheme *types.TransliterationScheme) (*Aksharamala, error) {
	if scheme == nil {
		return nil, fmt.Errorf("scheme cannot be nil")
	}

	context := NewContext()
	t := &Aksharamala{
		scheme:       scheme,
		buffer:       NewBuffer(),
		context:      context,
		categoryMaps: make(map[string]map[string][]string),
	}

	// Initialize virama handler
	t.viramaHandler = NewViramaHandler(scheme.Metadata.Virama, context)

	if err := t.buildMaps(); err != nil {
		return nil, fmt.Errorf("failed to build category maps: %w", err)
	}

	return t, nil
}

// buildMaps preprocesses the scheme mappings for efficient lookup
func (a *Aksharamala) buildMaps() error {
	for category, section := range a.scheme.Categories {
		catMap := make(map[string][]string)
		for _, mapping := range section.Mappings {
			if len(mapping.LHS) == 0 || len(mapping.RHS) == 0 {
				return fmt.Errorf("invalid mapping in category %s: empty LHS or RHS", category)
			}
			for _, lhs := range mapping.LHS {
				catMap[lhs] = mapping.RHS
			}
		}
		a.categoryMaps[category] = catMap
	}
	return nil
}

// Transliterate converts a string using the loaded scheme
func (a *Aksharamala) Transliterate(input string) (string, error) {
	var result strings.Builder

	// Clear buffer and reset state
	a.buffer = NewBuffer()

	// Add all characters to buffer first
	for _, r := range input {
		a.buffer.Append(r)
	}

	// Process until buffer is empty
	for a.buffer.Len() > 0 {
		// Get longest match from current position
		match, outputs := a.findLongestMatch()

		if match == "" {
			// No match found, consume one character and continue
			result.WriteRune(a.buffer.First())
			a.buffer.RemoveFirst()
			continue
		}

		matchLen := len(match)

		// Handle vowel after consonant
		if len(outputs) > 1 && isConsonant(a.context.LastChar) {
			if outputs[1] == "" {
				// Skip implicit 'a'
				a.buffer.Remove(matchLen)
				continue
			}
			// Add matra
			result.WriteString(outputs[1])
			a.buffer.Remove(matchLen)
			continue
		}

		// Handle consonants
		if len(outputs[0]) > 0 && isConsonant([]rune(outputs[0])[0]) {
			output := outputs[0]
			currentChar := []rune(output)[0]

			// In the consonant handling section of Transliterate:
			// Remove current match from buffer to look ahead
			remainingBuffer := a.buffer.String()[matchLen:]
			a.buffer.Remove(matchLen)

			needsVirama := false

			// Look ahead - we need virama if next character is a consonant
			if len(remainingBuffer) > 0 {
				// Get next match
				nextMatch, nextOutputs := a.findLongestMatch()

				if nextMatch != "" {
					// Check if next is a consonant (no vowel form) or is followed by a consonant
					isNextConsonant := len(nextOutputs) == 1 ||
						(len(nextOutputs) > 1 && nextOutputs[1] == "")

					// Don't add virama if next char is a vowel
					if isNextConsonant && !strings.ContainsAny(nextMatch, "aeiou") {
						needsVirama = true
					}
				}

				// Restore buffer
				a.buffer = NewBuffer()
				for _, ch := range remainingBuffer {
					a.buffer.Append(ch)
				}
			}

			// Update context
			a.context.LastChar = currentChar

			// Write output
			result.WriteString(output)
			if needsVirama {
				result.WriteString("्")
			}
			continue
		}

		// Default case
		result.WriteString(outputs[0])
		if len(outputs[0]) > 0 {
			a.context.LastChar = []rune(outputs[0])[0]
		}
		a.buffer.Remove(matchLen)
	}

	return result.String(), nil
}

// findLongestMatch finds the longest matching sequence in the current buffer
func (a *Aksharamala) findLongestMatch() (string, []string) {
	bufStr := a.buffer.String()

	// Try each category in priority order
	categories := []string{"consonants", "vowels", "others", "digits"}
	for _, category := range categories {
		mappings, exists := a.categoryMaps[category]
		if !exists {
			continue
		}

		// Try progressively shorter substrings
		for i := len(bufStr); i > 0; i-- {
			if outputs, ok := mappings[bufStr[:i]]; ok {
				return bufStr[:i], outputs
			}
		}
	}

	return "", nil
}
