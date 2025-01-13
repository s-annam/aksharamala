package translit

import "aks.go/internal/types"

// LookupResult represents the result of a lookup operation.
type LookupResult struct {
	Output   string
	Category string
}

// Precomputed lookup table for efficient transliteration.
type LookupTable map[string]LookupResult

// BuildLookupTable constructs a precomputed lookup table from a transliteration scheme.
func BuildLookupTable(scheme *types.TransliterationScheme) LookupTable {
	table := make(LookupTable)
	for category, section := range scheme.Categories {
		for _, mapping := range section.Mappings {
			for _, lhs := range mapping.LHS {
				table[lhs] = LookupResult{
					Output:   mapping.RHS[0],
					Category: category,
				}
			}
		}
	}
	return table
}

// Lookup performs a fast lookup using the precomputed table.
func (table LookupTable) Lookup(char string) LookupResult {
	if result, exists := table[char]; exists {
		return result
	}
	return LookupResult{Output: "", Category: "other"}
}
