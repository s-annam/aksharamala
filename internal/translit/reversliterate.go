package translit

import (
	"strings"

	"aks.go/internal/core"
	"aks.go/internal/types"
)

func (a *Aksharamala) lookupForReversliteration(str string) core.LookupResult {
	// First try conjuncts section for the full string
	if section, exists := a.activeScheme.Categories["conjuncts"]; exists {
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
						Category:  "conjuncts",
						Found:     true,
					}
				}
			}
		}
	}

	// If no conjunct match, do regular single character lookup
	categories := []string{"consonants", "others", "vowels", "matras", "digits"}
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
	_, viramaMode, err := types.ParseVirama(a.activeScheme.Metadata.Virama)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	runes := []rune(input)
	var lastResult core.LookupResult

	for i := 0; i < len(runes); {
		// Try to find a conjunct first - start with the longest possible sequence
		maxLen := len(runes) - i
		foundConjunct := false

		for length := maxLen; length > 1; length-- {
			candidateStr := string(runes[i : i+length])
			lookup := a.lookupForReversliteration(candidateStr)
			if lookup.Found && lookup.Category == "conjuncts" {
				if viramaMode == types.SmartMode {
					result.WriteString(lookup.Output + "a")
				} else {
					result.WriteString(lookup.Output)
				}
				lastResult = lookup
				i += length
				foundConjunct = true
				break
			}
		}

		if !foundConjunct {
			// Single character processing
			lookup := a.lookupForReversliteration(string(runes[i]))
			if !lookup.Found {
				result.WriteString(string(runes[i]))
				i++
				continue
			}

			if lookup.Category == "consonants" {
				if lastResult.Category == "consonants" && viramaMode == types.SmartMode && lookup.AltOutput != "" {
					result.WriteString(lookup.AltOutput)
				} else if viramaMode == types.SmartMode {
					result.WriteString(lookup.Output + "a")
				} else {
					result.WriteString(lookup.Output)
				}
			} else {
				result.WriteString(lookup.Output)
			}
			lastResult = lookup
			i++
		}
	}

	return result.String(), nil
}
