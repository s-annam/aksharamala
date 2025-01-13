// Updated aksharamala.go for Stage 2 with context management
package translit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TransliterationScheme struct {
	Categories map[string][]Mapping `json:"categories"`
}

type Mapping struct {
	LHS []string `json:"lhs"`
	RHS []string `json:"rhs"`
}

type Context struct {
	LastCharCategory string // Tracks the category of the last processed character (e.g., consonant, vowel, other)
	LastOutput       string // Tracks the last output character
}

func NewContext() *Context {
	return &Context{
		LastCharCategory: "",
		LastOutput:       "",
	}
}

type Aksharamala struct {
	scheme  *TransliterationScheme
	context *Context
}

// NewAksharamala initializes a new Aksharamala instance.
func NewAksharamala(schemePath string) (*Aksharamala, error) {
	data, err := os.ReadFile(schemePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scheme: %w", err)
	}

	scheme := &TransliterationScheme{}
	if err := json.Unmarshal(data, scheme); err != nil {
		return nil, fmt.Errorf("failed to parse scheme: %w", err)
	}

	return &Aksharamala{
		scheme:  scheme,
		context: NewContext(),
	}, nil
}

// Transliterate performs transliteration for consonants, vowels, and mixed input.
func (a *Aksharamala) Transliterate(input string) string {
	var result strings.Builder
	for _, char := range input {
		output := a.lookup(string(char))
		if output != "" {
			// Update context based on the output
			category := a.getCategory(output)
			a.context.LastCharCategory = category
			a.context.LastOutput = output
			result.WriteString(output)
		} else {
			// For unmatched characters, treat as "other"
			result.WriteString(string(char))
			// Reset context
			a.context.LastCharCategory = "other"
			a.context.LastOutput = string(char)
		}
	}
	return result.String()
}

// lookup finds the transliteration for a single character.
func (a *Aksharamala) lookup(char string) string {
	for category, mappings := range a.scheme.Categories {
		for _, mapping := range mappings {
			for _, lhs := range mapping.LHS {
				if lhs == char {
					// Use matra (RHS[1]) if the previous character is a consonant
					if category == "vowels" && a.context.LastCharCategory == "consonants" && len(mapping.RHS) > 1 {
						return mapping.RHS[1]
					}
					return mapping.RHS[0] // Use full form otherwise
				}
			}
		}
	}
	return "" // No match found
}

// getCategory determines the category of the output character.
func (a *Aksharamala) getCategory(output string) string {
	for category, mappings := range a.scheme.Categories {
		for _, mapping := range mappings {
			for _, rhs := range mapping.RHS {
				if rhs == output {
					return category
				}
			}
		}
	}
	return "other"
}
