// This file is part of Aksharamala (aks.go).
//
// Aksharamala is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Aksharamala is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Aksharamala. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"aks.go/internal/types"
)

const (
	defaultVersion = "2025.1"            // Version for the generated file
	defaultLicense = "AGPL-3.0-or-later" // License for the generated file
)

// Validate mandatory fields
func validateMandatoryFields(scheme *types.TransliterationScheme) error {
	missingFields := []string{}

	if scheme.ID == "" {
		missingFields = append(missingFields, "id")
		scheme.ID = "unknown_id" // Assign default if required
	}
	if scheme.Name == "" {
		missingFields = append(missingFields, "name")
		scheme.Name = "Unnamed Transliteration"
	}
	if scheme.Language == "" {
		missingFields = append(missingFields, "language")
		scheme.Language = "unknown_language"
	}
	if scheme.Scheme == "" {
		missingFields = append(missingFields, "scheme")
		scheme.Scheme = "unknown_scheme"
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("mandatory fields missing: %s", strings.Join(missingFields, ", "))
	}
	return nil
}

// Convert 0x-prefixed hexadecimal Unicode values to their corresponding characters
func convertUnicode(value string) string {
	// Recognize optional context markers around the Unicode values
	match := regexp.MustCompile(`^((?:\([^)]*\))*)([^()]+)((?:\([^)]*\))*)?$`).FindStringSubmatch(value)
	var before, middle, after string

	if len(match) > 1 {
		before = match[1]
	}
	if len(match) > 2 {
		middle = match[2]
	}
	if len(match) > 3 {
		after = match[3]
	}

	// Split by commas to handle multiple values
	parts := strings.Split(middle, ",")
	var result strings.Builder

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "0x") {
			// Parse the hexadecimal value
			codePoint, err := strconv.ParseInt(part[2:], 16, 32)
			if err == nil {
				result.WriteRune(rune(codePoint)) // Append the Unicode character
			}
		} else {
			// Append non-Unicode parts as-is
			result.WriteString(part)
		}
	}

	// Combine the parts with the context markers
	return before + result.String() + after
}

func handleMappingMatch(part string) []string {
	splits := strings.Fields(part)
	for i, split := range splits {
		splits[i] = handleSingleMapping(split)
	}
	return splits
}

func handleSingleMapping(value string) string {
	// Handle left and right square brackets
	value = strings.ReplaceAll(value, "\\[", "\\u005C\\u005B")
	value = strings.ReplaceAll(value, "\\]", "\\u005C\\u005D")

	// Now, handled unescaped left and right square brackets. These are used to
	// represent context in the AKT based keymaps. Going forward, we will use
	// parentheses to represent context in the JSON based keymaps.
	value = strings.ReplaceAll(value, "[", "(")
	value = strings.ReplaceAll(value, "]", ")")

	// Handle left and right curly braces
	value = strings.ReplaceAll(value, "{", "\\u007B")
	value = strings.ReplaceAll(value, "}", "\\u007D")
	value = strings.ReplaceAll(value, "\\{", "\\u005C\\u007B")
	value = strings.ReplaceAll(value, "\\}", "\\u005C\\u007D")

	// Handle escaped backslashes
	value = strings.ReplaceAll(value, "\\\\", "\\u005C")

	// Handle old-style Unicode strings
	return convertUnicode(value)
}

func main() {
	// Parse command-line arguments
	inputFile, outputFile := parseArgs()

	// Read and parse input file
	scheme, err := readAndParseInput(inputFile)
	if err != nil {
		fmt.Printf("Error processing input file: %v\n", err)
		return
	}

	// Convert to compact scheme
	compactScheme, err := convertToCompactScheme(scheme, inputFile)
	if err != nil {
		fmt.Printf("Error converting to compact scheme: %v\n", err)
		return
	}

	// Format and write output
	if err := writeOutput(compactScheme, outputFile); err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		return
	}

	fmt.Printf("AKT file '%s' converted to JSON successfully: %s\n", inputFile, outputFile)
}

func parseArgs() (string, string) {
	defaultInput := "../../examples/example.akt"
	defaultOutput := "../../examples/example.aktj"

	// Command-line arguments
	flag.Parse()
	args := flag.Args()

	inputFile := defaultInput
	if len(args) > 0 {
		inputFile = args[0]
	}

	outputFile := defaultOutput
	if len(args) > 1 {
		outputFile = args[1]
	}
	return inputFile, outputFile
}

func readAndParseInput(inputFile string) (types.TransliterationScheme, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return types.TransliterationScheme{}, fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	return ParseAKTFile(file)
}

func convertToCompactScheme(scheme types.TransliterationScheme, inputFile string) (types.CompactTransliterationScheme, error) {
	// Override comments
	sourceFile := filepath.Base(inputFile)
	scheme.Comments = []string{
		fmt.Sprintf("Converted from %s.", sourceFile),
		"Distributed under the GNU Affero General Public License (AGPL).",
	}

	return types.ToCompactTransliterationScheme(scheme)
}

func writeOutput(scheme types.CompactTransliterationScheme, outputFile string) error {
	// Format the JSON
	formattedJSON, err := FormatSchemeJSON(scheme)
	if err != nil {
		return fmt.Errorf("error formatting JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(outputFile, []byte(formattedJSON), 0o644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
