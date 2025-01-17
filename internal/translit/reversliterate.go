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
	virama, viramaMode, err := types.ParseVirama(a.activeScheme.Metadata.Virama)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	runes := []rune(input)

	for i := 0; i < len(runes); {
		// 1. Try full conjunct matches first
		fullStr := string(runes[i:])
		lookup := a.lookupForReversliteration(fullStr)
		if lookup.Found && lookup.Category == "conjuncts" {
			if viramaMode == types.NormalMode {
				result.WriteString(lookup.Output + virama)
			} else {
				result.WriteString(lookup.Output)
			}
			i += utf8.RuneCountInString(fullStr) // Skip the full conjunct
			continue
		}

		// 2. Single character handling
		lookup = a.lookupForReversliteration(string(runes[i]))
		if !lookup.Found {
			result.WriteString(string(runes[i]))
			i++
			continue
		}

		// 3. Check if the next charcter is a matra
		matraUpNext := false
		if i+1 < len(runes) && a.lookupForReversliteration(string(runes[i+1])).Category == "matras" {
			matraUpNext = true
		}

		switch lookup.Category {
		case "consonants":
			result.WriteString(lookup.Output)
			if viramaMode == types.NormalMode {
				if !matraUpNext {
					result.WriteString(virama)
				}
			}
		case "matras":
			if lookup.Output != "\u0000" { // Ignore empty matra
				result.WriteString(lookup.Output)
			}
		case "vowels", "others", "digits":
			result.WriteString(lookup.Output)
		}

		i++
	}

	return result.String(), nil
}
