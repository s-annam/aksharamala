package translit

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"aks.go/internal/types"
)

// Aksharamala handles Indic script transliteration using configurable schemes
type Aksharamala struct {
	scheme       *types.TransliterationScheme
	buffer       *Buffer
	context      *Context
	categoryMaps map[string]map[string][]string
}

// NewAksharamala creates a new transliterator for the given scheme
func NewAksharamala(scheme *types.TransliterationScheme) (*Aksharamala, error) {
	if scheme == nil {
		return nil, fmt.Errorf("scheme cannot be nil")
	}

	t := &Aksharamala{
		scheme:       scheme,
		buffer:       NewBuffer(),
		context:      NewContext(),
		categoryMaps: make(map[string]map[string][]string),
	}

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

	for _, r := range input {
		tr := a.transliterateRune(r)
		if tr.Error != nil {
			return "", tr.Error
		}

		// Handle backspaces if needed
		if tr.BackspaceCount > 0 {
			current := result.String()
			if len(current) < tr.BackspaceCount {
				return "", fmt.Errorf("invalid backspace count")
			}
			result.Reset()
			result.WriteString(current[:len(current)-tr.BackspaceCount])
		}

		result.WriteString(tr.Output)
	}

	// Process any remaining buffer content
	if a.buffer.Len() > 0 {
		tr := a.flushBuffer()
		if tr.Error != nil {
			return "", tr.Error
		}
		result.WriteString(tr.Output)
	}

	return result.String(), nil
}

// transliterateRune handles single character transliteration
func (a *Aksharamala) transliterateRune(r rune) Result {
	// Handle English mode toggle
	if r == '#' {
		a.context.ToggleEnglishMode()
		return Result{Output: ""}
	}

	// In English mode, pass through unchanged
	if a.context.InEnglishMode {
		return Result{Output: string(r)}
	}

	// Add to buffer and try to match
	a.buffer.Append(r)

	match, output := a.findLongestMatch()
	if match == "" {
		// No match found, return first character if buffer has multiple chars
		if a.buffer.Len() > 1 {
			first := a.buffer.First()
			a.buffer.RemoveFirst()
			return Result{Output: string(first)}
		}
		return Result{Output: ""}
	}

	// Process the match
	matchLen := utf8.RuneCountInString(match)
	backspaces := matchLen - 1
	a.buffer.Remove(matchLen)

	// Apply any context-specific rules
	output = a.applyContextRules(output)

	return Result{
		Output:         output,
		BackspaceCount: backspaces,
	}
}

// findLongestMatch finds the longest matching sequence in the current buffer
func (a *Aksharamala) findLongestMatch() (string, string) {
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
				return bufStr[:i], outputs[0]
			}
		}
	}

	return "", ""
}

// applyContextRules handles context-specific modifications
func (a *Aksharamala) applyContextRules(output string) string {
	// Check for virama handling based on scheme metadata
	if vrm := a.scheme.Metadata.Virama; vrm != "" {
		output = a.handleVirama(output, vrm)
	}

	// Update context based on output
	if output != "" {
		a.context.UpdateWithOutput(output)
	}

	return output
}

// handleVirama applies virama rules based on scheme configuration
func (a *Aksharamala) handleVirama(output, virama string) string {
	// TODO: Implement virama handling logic based on scheme settings
	return output
}

// flushBuffer processes any remaining content in the buffer
func (a *Aksharamala) flushBuffer() Result {
	if a.buffer.Len() == 0 {
		return Result{}
	}

	// Return remaining content as-is
	content := a.buffer.String()
	a.buffer.Clear()
	return Result{Output: content}
}
