package translit

import (
	"fmt"
	"strings"

	"aks.go/internal/core"
	"aks.go/internal/types"
)

// Transliterate performs transliteration for the input string using the active scheme.
// Returns the transliterated string.
func (a *Aksharamala) Transliterate(input string) (string, error) {
	// Reset context for a clean state
	a.context = types.NewContext()
	a.context.Input = input
	a.viramaHandler = types.NewViramaHandler(a.viramaHandler.Mode, a.viramaHandler.Virama, a.context)

	var result strings.Builder
	length := len(input)
	for i := 0; i < length; {
		a.context.Position = i
		foundMatch := false

		// Handle space character
		if i < length && input[i] == ' ' {
			shouldAddVirama, shouldAddSpace := a.viramaHandler.HandleSpace()
			if shouldAddVirama {
				result.WriteString(a.viramaHandler.Virama)
			}
			if shouldAddSpace {
				result.WriteRune(' ')
			}
			i++
			a.context.LatestLookup = core.LookupResult{Output: " ", Category: "other"}
			continue
		}

		// Greedily match substrings to find the longest match
		for j := length - i; j > 0; j-- {
			if i+j <= length {
				combination := input[i : i+j]
				lookupResult := a.lookup(combination)

				if lookupResult.Output != "" {
					if lookupResult.Output == "\x00" && a.context.LatestLookup.Category == "consonants" {
						a.context.LatestLookup = lookupResult
						i += j // Move the index forward by the length of the match
						foundMatch = true
						break
					}

					// Parse and apply contextual rules
					baseOutput, rules := types.ParseContextualRules(lookupResult.Output)
					lookupResult.Output = baseOutput

					// Only add virama for regular consonants, not for word boundary markers
					if lookupResult.Category != "word_boundary" {
						nextCategory := a.getCategoryForRHS(lookupResult.Output)
						if a.viramaHandler.ShouldInsertVirama(lookupResult.Output, nextCategory) {
							result.WriteString(a.viramaHandler.Virama)
						}
					}

					result.WriteString(lookupResult.Output)

					// Apply any contextual rules
					if err := a.context.ApplyContextualRules(rules, &result); err != nil {
						// Log error but continue with transliteration
						fmt.Printf("Error applying contextual rules: %v\n", err)
					}

					a.context.LatestLookup = lookupResult
					i += j // Move the index forward by the length of the match
					foundMatch = true
					break // Exit the loop once a match is found
				}
			}
		}

		// If no match was found, process the current character
		if !foundMatch {
			char := string(input[i])
			lookupResult := a.lookup(char)
			if lookupResult.Output != "" {
				// Parse and apply contextual rules
				baseOutput, rules := types.ParseContextualRules(lookupResult.Output)
				lookupResult.Output = baseOutput

				// Only add virama for regular consonants, not for word boundary markers
				if lookupResult.Category != "word_boundary" {
					nextCategory := a.getCategoryForRHS(lookupResult.Output)
					if a.viramaHandler.ShouldInsertVirama(lookupResult.Output, nextCategory) {
						result.WriteString(a.viramaHandler.Virama)
					}
				}

				result.WriteString(lookupResult.Output)

				// Apply any contextual rules
				if err := a.context.ApplyContextualRules(rules, &result); err != nil {
					// Log error but continue with transliteration
					fmt.Printf("Error applying contextual rules: %v\n", err)
				}

				a.context.LatestLookup = lookupResult
			} else {
				result.WriteString(char)
				a.context.LatestLookup = core.LookupResult{Output: char, Category: "other"}
			}
			i++ // Move to the next character
		}
	}

	// Handle end of input
	if a.viramaHandler.HandleEndOfInput() {
		result.WriteString(a.viramaHandler.Virama)
	}

	return result.String(), nil
}

// lookup finds the transliteration for the given string.
// Returns the LookupResult for the character.
func (a *Aksharamala) lookup(combination string) core.LookupResult {
	for category, section := range a.activeScheme.Categories {
		for _, mapping := range section.Mappings.All() {
			for _, lhs := range mapping.LHS {
				if lhs == combination {
					rhs := mapping.RHS
					if len(rhs) == 0 {
						continue
					}

					matchLen := len([]rune(combination))

					// Check for word boundary variants first
					if len(rhs) > 1 {
						// Check if the second option has a word boundary condition
						if strings.Contains(rhs[1], "(W)") {
							if a.context.IsSeparator(matchLen) {
								// Remove the (W) marker and return the rest
								output := strings.Replace(rhs[1], "(W)", "", 1)
								// Mark this as a special category so virama isn't added
								return core.LookupResult{
									Output:      output,
									Category:    "word_boundary",
									Found:       true,
									MatchLength: matchLen,
								}
							}
						} else if category == "vowels" && a.context.LatestLookup.Category == "consonants" {
							// Use matra if the previous character is a consonant
							return core.LookupResult{
								Output:      rhs[1],
								Category:    category,
								Found:       true,
								MatchLength: matchLen,
							}
						}
					}

					// Use first option as default
					return core.LookupResult{
						Output:      rhs[0],
						Category:    category,
						Found:       true,
						MatchLength: matchLen,
					}
				}
			}
		}
	}

	// No match found
	return core.LookupResult{
		Output:      "",
		Category:    "other",
		Found:       false,
		MatchLength: 0,
	}
}

// getCategoryForRHS determines which category a character belongs to
func (a *Aksharamala) getCategoryForRHS(output string) string {
	for category, section := range a.activeScheme.Categories {
		for _, mapping := range section.Mappings.All() {
			for _, rhs := range mapping.RHS {
				if rhs == output {
					return category
				}
			}
		}
	}
	return "other"
}
