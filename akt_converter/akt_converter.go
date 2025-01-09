package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Struct for JSON output
type Metadata struct {
	Virama       string `json:"virama,omitempty"`        // Omit if empty
	FontName     string `json:"font_name,omitempty"`     // Omit if empty
	FontSize     int    `json:"font_size,omitempty"`     // Omit if 0
	IconEnabled  string `json:"icon_enabled,omitempty"`  // Omit if empty
	IconDisabled string `json:"icon_disabled,omitempty"` // Omit if empty
}

type CategoryEntry struct {
	LHS     string   `json:"lhs"`
	RHS     []string `json:"rhs"`
	Comment string   `json:"comment,omitempty"` // Omit if no comment
}
type Section struct {
	Comments []string        `json:"comments,omitempty"` // Section-level comments
	Mappings []CategoryEntry `json:"mappings"`
}
type TransliterationScheme struct {
	Comments   []string           `json:"comments,omitempty"` // File-level comments
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Language   string             `json:"language"`
	Scheme     string             `json:"scheme"`
	Metadata   Metadata           `json:"metadata"`
	Categories map[string]Section `json:"categories"`
}
type Mapping struct {
	LHS     string   `json:"lhs"`
	RHS     []string `json:"rhs"`
	Context string   `json:"context,omitempty"`
}

// Validate mandatory fields
func validateMandatoryFields(scheme *TransliterationScheme) error {
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

// Main parsing function
func parseFile(file *os.File) (TransliterationScheme, error) {
	scanner := bufio.NewScanner(file)
	scheme := TransliterationScheme{
		Categories: make(map[string]Section),
	}

	sectionPattern := regexp.MustCompile(`^#(\w+)#`)
	metadataPattern := regexp.MustCompile(`#(\w+)\s*=\s*(.+)#?$`)

	// Temporary storage for file-level comments
	var fileComments []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Match file-level comments (// ...)
		if strings.HasPrefix(line, "//") {
			fileComments = append(fileComments, strings.TrimPrefix(line, "//"))
			continue
		}

		// Stop collecting comments as soon as a non-comment "#" line is found
		if strings.HasPrefix(line, "#") {
			// Check if it's metadata
			if match := metadataPattern.FindStringSubmatch(line); match != nil {
				parseMetadata(line, &scheme)
				continue
			}

			// Check if it's a section header
			if match := sectionPattern.FindStringSubmatch(line); match != nil {
				sectionName := match[1]
				section, err := parseSection(scanner, line)
				if err != nil {
					return scheme, err
				}
				scheme.Categories[sectionName] = section
				continue
			}
		}
	}

	// Validate mandatory fields
	if err := validateMandatoryFields(&scheme); err != nil {
		return scheme, err
	}

	// Assign file-level comments to the scheme
	scheme.Comments = fileComments
	return scheme, scanner.Err()
}

// Parse a section with comments and mappings
func parseSection(scanner *bufio.Scanner, sectionName string) (Section, error) {
	section := Section{}
	var pendingComment string
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*?)(?:\s+//\s*(.*))?$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect a new section or end of file
		if strings.HasPrefix(line, "#") && line != sectionName {
			return section, nil
		}

		// Handle section-level comments before mappings
		if strings.HasPrefix(line, "//") {
			trimmedComment := strings.TrimPrefix(line, "//")
			if len(section.Mappings) == 0 {
				// Assign to section-level comments if no mappings yet
				section.Comments = append(section.Comments, trimmedComment)
			} else {
				// Treat as pending mapping-level comment
				pendingComment = trimmedComment
			}
			continue
		}

		// Match mappings with inline comments
		match := mappingPattern.FindStringSubmatch(line)
		if match != nil {
			entry := CategoryEntry{
				LHS:     match[1],
				RHS:     strings.Fields(match[2]),
				Comment: match[3], // Inline comment
			}

			// Attach pending comment
			if pendingComment != "" {
				entry.Comment = pendingComment
				pendingComment = ""
			}

			section.Mappings = append(section.Mappings, entry)
		}
	}

	return section, scanner.Err()
}

func parseMetadata(line string, scheme *TransliterationScheme) {
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
	outputFile := "output.json"

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
