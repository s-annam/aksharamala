package translit

import (
	"strings"

	"aks.go/internal/core"
	"aks.go/internal/types"
)

func (a *Aksharamala) lookupForReversliteration(str string) core.LookupResult {
	// Order matters - check conjuncts first, then other categories
	categories := []string{"conjuncts", "consonants", "others", "vowels", "matras", "digits"}

	for _, category := range categories {
		if section, exists := a.activeScheme.Categories[category]; exists {
			for _, mapping := range section.Mappings.All() {
				for _, lhs := range mapping.LHS {
					if lhs == str {
						var altOutput string
						if len(mapping.RHS) > 1 {
							altOutput = mapping.RHS[1]
						}
						return core.LookupResult{
							Output:    mapping.RHS[0],
							AltOutput: altOutput,
							Category:  category,
							Found:     true,
						}
					}
				}
			}
		}
	}
	return core.LookupResult{Found: false, Category: "other"}
}

func (a *Aksharamala) Reversliterate(input string) (string, error) {
	virama, viramaMode, err := types.ParseVirama(a.activeScheme.Metadata.Virama)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	runes := []rune(input)
	var lastResult core.LookupResult

	for i := 0; i < len(runes); i++ {
		// Skip explicit virama
		if string(runes[i]) == virama {
			continue
		}

		lookup := a.lookupForReversliteration(string(runes[i]))
		if !lookup.Found {
			result.WriteString(string(runes[i]))
			continue
		}

		switch lookup.Category {
		case "consonants":
			if lastResult.Category == "consonants" && viramaMode == types.SmartMode && lookup.AltOutput != "" {
				result.WriteString(lookup.AltOutput)
			} else if viramaMode == types.SmartMode {
				result.WriteString(lookup.Output + "a")
			} else {
				result.WriteString(lookup.Output)
			}
		case "conjuncts":
			// Same logic as consonants
			if lastResult.Category == "consonants" && viramaMode == types.SmartMode && lookup.AltOutput != "" {
				result.WriteString(lookup.AltOutput)
			} else if viramaMode == types.SmartMode {
				result.WriteString(lookup.Output + "a")
			} else {
				result.WriteString(lookup.Output)
			}
		default:
			result.WriteString(lookup.Output)
		}

		lastResult = lookup
	}

	return result.String(), nil
}
