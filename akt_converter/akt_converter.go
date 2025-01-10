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
	LHS     []string `json:"lhs"` // Updated to a slice for multiple LHSs
	RHS     []string `json:"rhs"`
	Comment string   `json:"comment,omitempty"`
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

func parseMapping(line string, lastMapping *CategoryEntry) *CategoryEntry {
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*?)(?:\s+//\s*(.*))?$`)
	lhsOnlyPattern := regexp.MustCompile(`^(\S+)$`)

	// Match full mappings
	if match := mappingPattern.FindStringSubmatch(line); match != nil {
		return &CategoryEntry{
			LHS:     []string{match[1]},
			RHS:     strings.Fields(match[2]),
			Comment: match[3], // Inline comment
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

// Main parsing function
func parseFile(file *os.File) (TransliterationScheme, error) {
	scanner := bufio.NewScanner(file)
	scheme := TransliterationScheme{
		Categories: make(map[string]Section),
	}

	metadataPattern := regexp.MustCompile(`#(\w+)\s*=\s*(.+)#?$`)
	sectionPattern := regexp.MustCompile(`^#(\w+)#`)                                   // Regular sections
	pseudoSectionPattern := regexp.MustCompile(`^\/\/\s*=\*=\s*(\w+)(?:\s*=\*=\s*)?$`) // Pseudo-sections

	var currentCategory string
	var section Section
	var fileComments []string
	var lastMapping *CategoryEntry

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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
			section = Section{}
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
			section = Section{}
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
