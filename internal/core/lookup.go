package core

// LookupResult represents the result of a lookup operation.
type LookupResult struct {
	Output   string
	Category string
}

// Precomputed lookup table for efficient transliteration.
type LookupTable map[string]LookupResult

// Lookup performs a fast lookup using the precomputed table.
func (table LookupTable) Lookup(char string) LookupResult {
	if result, exists := table[char]; exists {
		return result
	}
	return LookupResult{Output: "", Category: "other"}
}
