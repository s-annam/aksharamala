package core

// LookupResult represents the result of a lookup operation.
type LookupResult struct {
	Output      string // Primary output (first RHS)
	AltOutput   string // Alternate output (second RHS)
	Category    string // Category of the mapping
	Found       bool   // Whether the lookup found a match
	MatchLength int    // Length of the current match in runes
}

// LookupTable maps input strings to their lookup results
type LookupTable map[string]LookupResult

// Lookup performs a fast lookup using the table
func (table LookupTable) Lookup(char string) LookupResult {
	if result, exists := table[char]; exists {
		return result
	}
	return LookupResult{Found: false, Category: "other"}
}

// Lookup maintains the same simple interface but uses direct mapping
func (m *Mappings) Lookup(lhs string) (string, bool) {
	// Create a simple lookup table just from this mapping's entries
	table := make(LookupTable)
	for _, entry := range m.entries {
		for _, lhsEntry := range entry.LHS {
			if len(entry.RHS) > 0 {
				table[lhsEntry] = LookupResult{
					Output: entry.RHS[0],
					Found:  true,
				}
			}
		}
	}

	result := table.Lookup(lhs)
	return result.Output, result.Found
}
