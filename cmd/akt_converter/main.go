package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"aks.go/internal/types"
)

const currentVersion = "2025.1" // Define the current version for the JSON format

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

// Parse mappings and handle multiple LHS
func parseMapping(line string, lastMapping *types.CategoryEntry) *types.CategoryEntry {
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*?)(?:\s+//\s*(.*))?$`)
	lhsOnlyPattern := regexp.MustCompile(`^(\S+)$`)

	// Match full mappings
	if match := mappingPattern.FindStringSubmatch(line); match != nil {
		return &types.CategoryEntry{
			LHS:     []string{match[1]},
			RHS:     strings.Fields(match[2]),
			Comment: match[3],
		}
	}

	// Match LHS-only lines and attach to the last mapping
	if match := lhsOnlyPattern.FindStringSubmatch(line); match != nil {
		if lastMapping != nil {
			lastMapping.LHS = append(lastMapping.LHS, match[1])
		}
		return nil
	}

	return nil
}

// Parse the AKT file into TransliterationScheme
func parseFile(file *os.File) (types.TransliterationScheme, error) {
	scanner := bufio.NewScanner(file)
	scheme := types.TransliterationScheme{
		Version:    currentVersion, // Assign the current version
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
	inputFile := "example.akt"
	outputFile := "output.aktj"

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scheme, err := parseFile(file)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(scheme); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Printf("AKT file converted to JSON successfully: %s\n", outputFile)
}
