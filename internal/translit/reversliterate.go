package translit

import (
	"strings"
	"unicode/utf8"

	"aks.go/internal/types"
)

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
		lookup := a.lookup(fullStr)
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
		lookup = a.lookup(string(runes[i]))
		if !lookup.Found {
			result.WriteString(string(runes[i]))
			i++
			continue
		}

		// 3. Check if the next charcter is a matra
		matraUpNext := false
		if i+1 < len(runes) && a.lookup(string(runes[i+1])).Category == "matras" {
			matraUpNext = true
		}

		switch lookup.Category {
		case "consonants":
			result.WriteString(lookup.Output)
			if !matraUpNext {
				result.WriteString(virama)
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
