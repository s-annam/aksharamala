package translit

import (
	"fmt"
	"strings"
	"unicode"
)

// ViramaMode represents different virama handling behaviors
type ViramaMode int

const (
	ViramaModeSmart ViramaMode = iota
	ViramaModeNormal
	ViramaModeDouble
	ViramaModeRepeat
)

// ViramaHandler handles virama-related operations during transliteration
type ViramaHandler struct {
	char    rune       // The virama character to use
	mode    ViramaMode // Current virama mode
	context *Context   // Reference to transliteration context
}

// NewViramaHandler creates a new ViramaHandler from metadata configuration
func NewViramaHandler(viramaStr string, ctx *Context) *ViramaHandler {
	char, mode := parseViramaConfig(viramaStr)
	return &ViramaHandler{
		char:    char,
		mode:    mode,
		context: ctx,
	}
}

// parseViramaConfig extracts virama character and mode from metadata
func parseViramaConfig(viramaStr string) (rune, ViramaMode) {
	if viramaStr == "" {
		return '्', ViramaModeSmart // Default to Devanagari virama in smart mode
	}

	parts := strings.Split(viramaStr, ",")
	if len(parts) != 2 {
		return '्', ViramaModeSmart
	}

	// Parse virama character (handling hex format "0x094D")
	var viramaChar rune
	if strings.HasPrefix(parts[0], "0x") {
		_, err := fmt.Sscanf(parts[0], "0x%X", &viramaChar)
		if err != nil {
			viramaChar = '्'
		}
	} else {
		viramaChar = []rune(strings.TrimSpace(parts[0]))[0]
	}

	// Parse mode
	mode := strings.ToLower(strings.TrimSpace(parts[1]))
	switch mode {
	case "normal":
		return viramaChar, ViramaModeNormal
	case "double":
		return viramaChar, ViramaModeDouble
	case "repeat":
		return viramaChar, ViramaModeRepeat
	default:
		return viramaChar, ViramaModeSmart
	}
}

// ApplyVirama handles virama insertion based on current mode and context
func (vh *ViramaHandler) ApplyVirama(output string) string {
	if output == "" || vh.context.InEnglishMode {
		return output
	}

	outputRunes := []rune(output)
	lastRune := outputRunes[len(outputRunes)-1]

	// Don't apply virama if we're handling a vowel mark
	if isVowelMark(lastRune) {
		return output
	}

	// Don't add virama for first character
	if vh.context.LastChar == 0 {
		return output
	}

	// Check if the current character is a consonant
	if !isConsonant(lastRune) {
		return output
	}

	prevChar := vh.context.LastChar

	switch vh.mode {
	case ViramaModeSmart:
		// Insert virama between consecutive consonants
		if isConsonant(prevChar) {
			return string(vh.char) + string(lastRune)
		}

	case ViramaModeNormal:
		// Always insert virama after consonant
		if isConsonant(prevChar) {
			return string(vh.char) + string(lastRune)
		}

	case ViramaModeDouble:
		// Insert virama for doubled consonants
		if prevChar == lastRune {
			return string(vh.char) + string(lastRune)
		}

	case ViramaModeRepeat:
		// Check if current consonant matches the previous one
		if isConsonant(prevChar) && prevChar == lastRune {
			return string(vh.char) + string(lastRune)
		}
	}

	return output
}

// isConsonant checks if a character is a consonant
// This is a simplified check - we should use proper Unicode properties
// or maintain a map of consonant characters based on the scheme
func isConsonant(r rune) bool {
	return unicode.Is(unicode.Devanagari, r) && !isVowelMark(r)
}

// isVowelMark checks if a character is a vowel mark (matra)
func isVowelMark(r rune) bool {
	// Devanagari vowel marks range
	return r >= '\u093A' && r <= '\u094F'
}
