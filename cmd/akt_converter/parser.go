package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"aks.go/internal/types"
)

// ParseAKTFile parses an AKT file into a TransliterationScheme
func ParseAKTFile(file *os.File) (types.TransliterationScheme, error) {
	scanner := bufio.NewScanner(file)
	scheme := types.TransliterationScheme{
		Version:    defaultVersion, // Assign the current version
		License:    defaultLicense, // Assign the default license
		Categories: make(map[string]types.Section),
	}

	metadataPattern := regexp.MustCompile(`#(\w+)\s*=\s*(.+)#?$`)
	sectionPattern := regexp.MustCompile(`^#(\w+)#`)                                  // Regular sections
	pseudoSectionPattern := regexp.MustCompile(`^\/\/\s*=*\*=*\s*(.+?)\s*=*\*=*\s*$`) // Pseudo-sections

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

			// Get or create the section
			currentCategory = match[1]
			section = getOrCreateSection(scheme, currentCategory)
			lastMapping = nil // Reset last mapping for the new section
			continue
		}

		// Match pseudo-sections
		if match := pseudoSectionPattern.FindStringSubmatch(line); match != nil {
			// Save the previous section
			if currentCategory != "" {
				scheme.Categories[currentCategory] = section
			}

			// Get or create the (pseudo) section, consider only the first word
			currentCategory = strings.ToLower(strings.Fields(match[1])[0])
			section = getOrCreateSection(scheme, currentCategory)
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

// Get an existing section or create a new one if it doesn't exist
func getOrCreateSection(scheme types.TransliterationScheme, currentCategory string) types.Section {
	if existingSection, exists := scheme.Categories[currentCategory]; exists {
		return existingSection // Use the existing section
	}
	return types.Section{} // Create a new section if it doesn't exist
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

func normalizeComment(comment string) string {
	// Trim spaces
	comment = strings.TrimSpace(comment)

	// Remove surrounding `=*=` if present
	if strings.HasPrefix(comment, "=*=") && strings.HasSuffix(comment, "=*=") {
		comment = strings.TrimPrefix(comment, "=*=")
		comment = strings.TrimSuffix(comment, "=*=")
	}

	return comment
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
			Comment: normalizeComment(match[3]),
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
