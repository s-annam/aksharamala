// Simplified aksharamala.go for Stage 1
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

type Aksharamala struct {
	scheme *TransliterationScheme
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

	return &Aksharamala{scheme: scheme}, nil
}

// Transliterate performs basic transliteration for consonants and vowels.
func (a *Aksharamala) Transliterate(input string) string {
	var result strings.Builder
	for _, char := range input {
		output := a.lookup(string(char))
		result.WriteString(output)
	}
	return result.String()
}

// lookup finds the transliteration for a single character.
func (a *Aksharamala) lookup(char string) string {
	for _, category := range a.scheme.Categories {
		for _, mapping := range category {
			for _, lhs := range mapping.LHS {
				if lhs == char {
					return mapping.RHS[0] // Use the full form for now
				}
			}
		}
	}
	return char // Default to input if no match
}
