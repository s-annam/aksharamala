package translit

import (
	"strings"

	"aks.go/internal/types"
)

func (a *Aksharamala) Reversliterate(input string) (string, error) {
	// Reset context for a clean state
	a.context = types.NewContext()
	a.context.Input = input

	virama, viramaMode, err := types.ParseVirama(a.activeScheme.Metadata.Virama)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	runes := []rune(input)
	length := len(runes)

	for i := 0; i < length; {
		a.context.Position = i
		foundMatch := false

		// Greedily match substrings to find the longest match
		for j := length - i; j > 0; j-- {
			if i+j <= length {
				substr := string(runes[i : i+j])
				lookup := a.lookup(substr)

				if lookup.Found {
					// Handle based on category
					switch lookup.Category {
					case "consonants":
						result.WriteString(lookup.Output)
						// Add virama if we're at the end OR if next char isn't a matra
						if i+j >= length ||
							(i+j < length && a.lookup(string(runes[i+j])).Category != "matras") {
							if viramaMode == types.NormalMode {
								result.WriteString(virama)
							} else if viramaMode == types.SmartMode && !a.context.IsSeparator() {
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

					i += j // Move the index forward by the length of the match
					foundMatch = true
					break
				}
			}
		}

		// If no match was found, copy the character as is
		if !foundMatch {
			result.WriteString(string(runes[i]))
			i++
		}
	}

	return result.String(), nil
}
