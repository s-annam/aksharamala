package translit

import (
	"aks.go/internal/core"
	"aks.go/internal/types"
)

// BuildLookupTable constructs a precomputed lookup table from a transliteration scheme
func BuildLookupTable(scheme *types.TransliterationScheme) core.LookupTable {
	table := make(core.LookupTable)
	for category, section := range scheme.Categories {
		for _, mapping := range section.Mappings.All() {
			for _, lhs := range mapping.LHS {
				result := core.LookupResult{
					Category: category,
					Found:    true,
				}

				if len(mapping.RHS) > 0 {
					result.Output = mapping.RHS[0]
					if len(mapping.RHS) > 1 {
						result.AltOutput = mapping.RHS[1]
					}
				}

				table[lhs] = result
			}
		}
	}
	return table
}
