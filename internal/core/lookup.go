package core

// LookupResult represents the result of a lookup operation.
// It contains the output string and the category of the character.
type LookupResult struct {
	Output   string // The transliterated output for the character
	Category string // The category of the character (e.g., consonant, vowel)
}

// Precomputed lookup table for efficient transliteration.
type LookupTable map[string]LookupResult

// Lookup performs a fast lookup using the precomputed table.
// It takes a character as input and returns the corresponding LookupResult.
// If the character is not found, it returns a LookupResult with an empty output
// and the category set to "other".
func (table LookupTable) Lookup(char string) LookupResult {
	if result, exists := table[char]; exists {
		return result
	}
	return LookupResult{Output: "", Category: "other"}
}
