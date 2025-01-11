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
	"bufio"
	"encoding/json"
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

// Parse mappings and handle multiple LHS
func parseMapping(line string, lastMapping *types.CategoryEntry) *types.CategoryEntry {
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*?)(?:\s+//\s*(.*))?$`)
	lhsOnlyPattern := regexp.MustCompile(`^(\S+)$`)

	// Match full mappings
	if match := mappingPattern.FindStringSubmatch(line); match != nil {
		return &types.CategoryEntry{
			LHS:     handleMappingMatch(match[1]),
			RHS:     handleMappingMatch(match[2]),
			Comment: match[3],
		}
	}

	// Match LHS-only lines and attach to the last mapping
	if match := lhsOnlyPattern.FindStringSubmatch(line); match != nil {
		if lastMapping != nil {
			// Append the new LHS to the last mapping's LHS
			x := handleMappingMatch(match[1])[0]
			lastMapping.LHS = append(lastMapping.LHS, x)
		}
		return nil
	}

	return nil
}

// Parse the AKT file into TransliterationScheme
func parseFile(file *os.File) (types.TransliterationScheme, error) {
	scanner := bufio.NewScanner(file)
	scheme := types.TransliterationScheme{
		Version:    defaultVersion, // Assign the current version
		License:    defaultLicense, // Assign the default license
		Categories: make(map[string]types.Section),
	}

	metadataPattern := regexp.MustCompile(`#(\w+)\s*=\s*(.+)#?$`)
	sectionPattern := regexp.MustCompile(`^#(\w+)#`)                                   // Regular sections
	pseudoSectionPattern := regexp.MustCompile(`^\/\/\s*=\*=\s*(\w+)(?:\s*=\*=\s*)?$`) // Pseudo-sections

	var currentCategory string
	var section types.Section
	var fileComments []string
	var lastMapping *types.CategoryEntry

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for end-of-file marker
		if strings.EqualFold(line, "#end") {
			break
		}

		// Match metadata
		if match := metadataPattern.FindStringSubmatch(line); match != nil {
			parseMetadata(line, &scheme)
			continue
		}

		// Match regular sections
		if match := sectionPattern.FindStringSubmatch(line); match != nil {
			// Save the previous section
			if currentCategory != "" {
				scheme.Categories[currentCategory] = section
			}

			// Start a new section
			currentCategory = match[1]
			section = types.Section{}
			lastMapping = nil // Reset last mapping for the new section
			continue
		}

		// Match pseudo-sections
		if match := pseudoSectionPattern.FindStringSubmatch(line); match != nil {
			// Save the previous section
			if currentCategory != "" {
				scheme.Categories[currentCategory] = section
			}

			// Start a new (pseudo) section
			currentCategory = strings.ToLower(match[1])
			section = types.Section{}
			lastMapping = nil // Reset last mapping for the new section
			continue
		}

		// Match file-level comments
		if strings.HasPrefix(line, "//") {
			fileComments = append(fileComments, strings.TrimPrefix(line, "//"))
			continue
		}

		// Match mappings within the current section
		entry := parseMapping(line, lastMapping)
		if entry != nil {
			section.Mappings = append(section.Mappings, *entry)
			lastMapping = &section.Mappings[len(section.Mappings)-1] // Track last mapping
		}
	}

	// Save the last section
	if currentCategory != "" {
		scheme.Categories[currentCategory] = section
	}

	// Save file-level comments
	scheme.Comments = fileComments

	if err := validateMandatoryFields(&scheme); err != nil {
		return scheme, err
	}

	return scheme, scanner.Err()
}

// Parse metadata fields from the AKT file
func parseMetadata(line string, scheme *types.TransliterationScheme) {
	metadataPattern := regexp.MustCompile(`#(\w+)\s*=\s*(.+)#?$`)
	match := metadataPattern.FindStringSubmatch(line)
	if match == nil {
		return
	}

	key := strings.ToLower(match[1])
	value := strings.TrimSpace(strings.TrimRight(match[2], "#"))

	switch key {
	case "id":
		scheme.ID = value
	case "name":
		scheme.Name = value
	case "language":
		scheme.Language = value
	case "scheme":
		scheme.Scheme = value
	case "virama":
		scheme.Metadata.Virama = value
	}
}

func main() {
	// Default paths
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

	// Open and parse the input AKT file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file '%s': %v\n", inputFile, err)
		return
	}
	defer file.Close()

	scheme, err := parseFile(file)
	if err != nil {
		fmt.Printf("Error parsing AKT file: %v\n", err)
		return
	}

	// Override comments
	sourceFile := filepath.Base(inputFile)
	scheme.Comments = []string{
		fmt.Sprintf("Keymap generated from %s. "+
			"Manually reviewed and refined for accuracy.", sourceFile),
		"Distributed under GNU Affero General Public License (AGPL) " +
			"as with the rest of the Aksharamala project.",
	}

	// Create the output JSON file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file '%s': %v\n", outputFile, err)
		return
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(scheme); err != nil {
		fmt.Printf("Error encoding JSON to output file: %v\n", err)
		return
	}

	fmt.Printf("AKT file '%s' converted to JSON successfully: %s\n", inputFile, outputFile)
}
