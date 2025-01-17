package translit

import (
	"strings"
	"unicode/utf8"

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
		// 1. Try full conjunct matches first
		fullStr := string(runes[i:])
		lookup := a.lookupForReversliteration(fullStr)
		if lookup.Found && lookup.Category == "conjuncts" {
			if viramaMode == types.SmartMode {
				result.WriteString(lookup.Output + "a")
			} else {
				result.WriteString(lookup.Output)
			}
			lastResult = lookup
			i += utf8.RuneCountInString(fullStr) // Skip the full conjunct
			continue
		}

		// 2. Look for geminated consonants
		if i+2 < len(runes) {
			currentChar := string(runes[i])
			currentLookup := a.lookupForReversliteration(currentChar)
			if currentLookup.Category == "consonants" {
				if string(runes[i+1]) == "्" {
					nextChar := string(runes[i+2])
					nextLookup := a.lookupForReversliteration(nextChar)
					if nextLookup.Found && nextLookup.Output == currentLookup.Output {
						if viramaMode == types.SmartMode {
							result.WriteString(currentLookup.Output + currentLookup.Output + "a")
						} else {
							result.WriteString(currentLookup.Output)
						}
						lastResult = currentLookup
						i += 3
						continue
					}
				}
			}
		}

		// 3. Single character handling
		lookup = a.lookupForReversliteration(string(runes[i]))
		if !lookup.Found {
			result.WriteString(string(runes[i]))
			i++
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
		case "matras":
			result.WriteString(lookup.Output)
		case "vowels", "others", "digits":
			result.WriteString(lookup.Output)
		}

		lastResult = lookup
		i++
	}

	return result.String(), nil
}
