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

	"aks.go/internal/core"
	"aks.go/internal/types"
	"aks.go/logger"
	"go.uber.org/zap"
	"encoding/json"
)

const (
	defaultVersion = "2025.1"            // Version for the generated file
	defaultLicense = "AGPL-3.0-or-later" // License for the generated file
)

// convertUnicode converts 0x-prefixed hexadecimal Unicode values to their corresponding characters.
// It takes a string value and returns the converted character as a string.
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

// handleMappingMatch processes a part of the input and returns the corresponding mappings.
// It takes a string part and returns a slice of strings.
func handleMappingMatch(part string) []string {
	splits := strings.Fields(part)
	for i, split := range splits {
		splits[i] = handleSingleMapping(split)
	}
	return splits
}

// handleSingleMapping processes a single mapping value and returns the corresponding character.
// It takes a string value and returns the mapped character as a string.
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

// main is the entry point of the application. It processes command-line flags,
// reads input files, and performs transliteration based on the provided keymaps.
// It returns an error if any step of the process fails.
func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "Enable debug logging")
	updateOnly := flag.Bool("update-only", false, "Update existing mappings in-place while keeping sections intact")
	noUpdate := flag.Bool("no-update", false, "Force creating a new file even if output exists")
	flag.Parse()

	// Initialize the logger
	logger.InitLogger(*debug)
	defer logger.Sync()

	// Parse input and output
	inputFile, outputFile := parseArgs()

	// Check if output file exists and determine update mode
	shouldUpdate := *updateOnly
	if !*noUpdate {
		if _, err := os.Stat(outputFile); err == nil {
			shouldUpdate = true
		}
	}

	logger.Info("Starting AKT conversion",
		zap.String("inputFile", inputFile),
		zap.String("outputFile", outputFile),
		zap.Bool("updateOnly", shouldUpdate))

	// Process the input file
	scheme, err := readAndParseInput(inputFile)
	if err != nil {
		logger.Error("Error reading input file", zap.String("inputFile", inputFile), zap.Error(err))
		return
	}

	// If in update mode, read the existing output file
	var existingScheme *types.TransliterationScheme
	if shouldUpdate {
		if existingFile, err := os.ReadFile(outputFile); err == nil {
			var existing types.TransliterationScheme
			if err := json.Unmarshal(existingFile, &existing); err == nil {
				existingScheme = &existing
				logger.Info("Successfully loaded existing scheme", zap.String("outputFile", outputFile))
			} else {
				logger.Warn("Failed to parse existing output file, will create new",
					zap.String("outputFile", outputFile),
					zap.Error(err))
			}
		} else {
			logger.Warn("Failed to read existing output file, will create new",
				zap.String("outputFile", outputFile),
				zap.Error(err))
		}
	}

	// Convert to compact scheme
	compactScheme, err := convertToCompactScheme(scheme, inputFile, existingScheme)
	if err != nil {
		logger.Error("Error converting to compact scheme", zap.Error(err))
		return
	}

	// Format and write output
	if err := writeOutput(compactScheme, outputFile); err != nil {
		logger.Error("Error writing output file", zap.String("outputFile", outputFile), zap.Error(err))
		return
	}

	logger.Info("AKT conversion completed successfully", zap.String("outputFile", outputFile))
}

// parseArgs parses command-line arguments and returns the input and output file paths.
// It returns the default input and output file paths if no arguments are provided.
func parseArgs() (string, string) {
	defaultInput := "../../examples/example.akt"
	defaultOutput := "../../examples/example.aksj"

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

// readAndParseInput reads and parses the input file, returning a TransliterationScheme.
// It returns an error if the file cannot be opened or parsed.
func readAndParseInput(inputFile string) (types.TransliterationScheme, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return types.TransliterationScheme{}, fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	return ParseAKTFile(file)
}

// convertToCompactScheme converts a TransliterationScheme to a CompactTransliterationScheme.
// It takes the scheme and input file path, returning the compact scheme and any error encountered.
// The function also overrides the comments in the scheme with information about the conversion.
func convertToCompactScheme(scheme types.TransliterationScheme, inputFile string, existingScheme *types.TransliterationScheme) (types.CompactTransliterationScheme, error) {
	// Override comments
	sourceFile := filepath.Base(inputFile)
	scheme.Comments = []string{
		fmt.Sprintf("Converted from %s.", sourceFile),
		"Distributed under the GNU Affero General Public License (AGPL).",
	}

	if existingScheme != nil {
		// In update-only mode, start with existing scheme
		mergedScheme := types.TransliterationScheme{
			Version:    scheme.Version,
			License:    scheme.License,
			ID:         scheme.ID,
			Name:       scheme.Name,
			Language:   scheme.Language,
			Scheme:     scheme.Scheme,
			Metadata:   scheme.Metadata,
			Comments:   scheme.Comments,
			Categories: make(map[string]types.Section),
		}

		// First copy all existing categories to preserve structure
		for name, section := range existingScheme.Categories {
			mergedScheme.Categories[name] = section
		}

		// Process each mapping from input scheme
		for inputSection, inputContent := range scheme.Categories {
			for _, mapping := range inputContent.Mappings.All() {
				// Try to find this mapping in the existing scheme
				if section, idx, found := existingScheme.FindMapping(mapping.LHS); found {
					// Update mapping in its current section
					existingSection := mergedScheme.Categories[section]
					existingMappings := existingSection.Mappings.All()
					// Preserve original LHS while updating RHS and Comment
					existingMappings[idx].RHS = mapping.RHS
					existingMappings[idx].Comment = mapping.Comment
					existingSection.Mappings = core.NewMappings(existingMappings)
					mergedScheme.Categories[section] = existingSection
					logger.Info("Updated existing mapping",
						zap.Strings("lhs", existingMappings[idx].LHS),
						zap.String("in_section", section))
				} else {
					// Mapping not found, add to appropriate section from input
					if section, exists := mergedScheme.Categories[inputSection]; exists {
						mappings := section.Mappings.All()
						mappings = append(mappings, mapping)
						section.Mappings = core.NewMappings(mappings)
						mergedScheme.Categories[inputSection] = section
					} else {
						// Create new section if it doesn't exist
						mergedScheme.Categories[inputSection] = types.Section{
							Mappings: core.NewMappings([]core.Mapping{mapping}),
						}
					}
					logger.Info("Added new mapping",
						zap.Strings("lhs", mapping.LHS),
						zap.String("to_section", inputSection))
				}
			}
		}

		return types.ToCompactTransliterationScheme(mergedScheme)
	}

	return types.ToCompactTransliterationScheme(scheme)
}

// writeOutput writes the CompactTransliterationScheme to the specified output file.
// It takes the scheme and output file path, returning any error encountered.
// The function formats the scheme as JSON before writing it to the file.
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
