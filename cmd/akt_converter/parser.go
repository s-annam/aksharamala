package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"aks.go/internal/core"
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
	var lastMapping *core.Mapping

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

			// Start a new section or reuse an existing one
			currentCategory = match[1]
			section = *types.GetOrCreate(scheme, currentCategory)
			lastMapping = nil
			continue
		}

		// Match pseudo-sections
		if match := pseudoSectionPattern.FindStringSubmatch(line); match != nil {
			// Save the previous section
			if currentCategory != "" {
				scheme.Categories[currentCategory] = section
			}

			// Start a new pseudo-section or reuse an existing one
			currentCategory = strings.ToLower(strings.Fields(match[1])[0])
			section = *types.GetOrCreate(scheme, currentCategory)
			lastMapping = nil
			continue
		}

		// Match mappings
		entry := parseAndAddMapping(line, &section, lastMapping)
		if entry != nil {
			lastMapping = entry
		}
	}

	// Save the last section
	if currentCategory != "" {
		scheme.Categories[currentCategory] = section
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

// Normalize comments
func normalizeComment(comment string) string {
	comment = strings.TrimSpace(comment)
	if strings.HasPrefix(comment, "=*=") && strings.HasSuffix(comment, "=*=") {
		comment = strings.TrimPrefix(comment, "=*=")
		comment = strings.TrimSuffix(comment, "=*=")
	}
	return comment
}

// parseAndAddMapping parses a single line from the AKT file and updates the given section.
//
// Parameters:
// - line: The line to parse.
// - section: The section to which the parsed mapping should be added or updated.
// - lastMapping: The most recent mapping added, used to update LHS-only patterns.
//
// Returns:
// - The new mapping if a full mapping is found.
func parseAndAddMapping(line string, section *types.Section, lastMapping *core.Mapping) *core.Mapping {
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*?)(?:\s+//\s*(.*))?$`)
	lhsOnlyPattern := regexp.MustCompile(`^(\S+)$`)

	// Match full mappings
	if match := mappingPattern.FindStringSubmatch(line); match != nil {
		entry := &core.Mapping{
			LHS:     handleMappingMatch(match[1]),
			RHS:     handleMappingMatch(match[2]),
			Comment: normalizeComment(match[3]),
		}
		section.AddMapping(entry.LHS, entry.RHS, entry.Comment)
		return entry
	}

	// Match LHS-only lines
	if match := lhsOnlyPattern.FindStringSubmatch(line); match != nil {
		if lastMapping != nil {
			section.AppendLHSToMapping(lastMapping, match[1])
			return nil
		}
	}

	return nil
}
