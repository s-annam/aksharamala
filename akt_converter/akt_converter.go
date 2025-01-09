package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Struct for JSON output
type TransliterationScheme struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Language string    `json:"language"`
	Scheme   string    `json:"scheme"`
	Metadata Metadata  `json:"metadata"`
	Mappings []Mapping `json:"mappings"`
	Comments []string  `json:"comments"`
}

type Metadata struct {
	Virama       string `json:"virama"`
	FontName     string `json:"font_name"`
	FontSize     int    `json:"font_size"`
	IconEnabled  string `json:"icon_enabled"`
	IconDisabled string `json:"icon_disabled"`
}

type Mapping struct {
	LHS     string   `json:"lhs"`
	RHS     []string `json:"rhs"`
	Context string   `json:"context,omitempty"`
}

func main() {
	// Input and output files
	inputFile := "example.akt"  // Replace with your AKT file path
	outputFile := "output.json" // Output JSON file

	// Open the AKT file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Initialize struct to hold data
	var scheme TransliterationScheme
	scheme.Mappings = []Mapping{}
	scheme.Comments = []string{}

	// Regex patterns
	metadataPattern := regexp.MustCompile(`#([a-zA-Z_]+)\s*=\s*(.*)`)
	mappingPattern := regexp.MustCompile(`^(\S+)\s+(\S.*)$`) // LHS and RHS
	commentPattern := regexp.MustCompile(`^//(.*)$`)

	// Read the file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Match comments
		if match := commentPattern.FindStringSubmatch(line); match != nil {
			scheme.Comments = append(scheme.Comments, strings.TrimSpace(match[1]))
			continue
		}

		// Match metadata
		if match := metadataPattern.FindStringSubmatch(line); match != nil {
			key := strings.ToLower(match[1])
			value := strings.TrimSpace(match[2])
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
			case "font_name":
				scheme.Metadata.FontName = value
			case "font_size":
				// Parse font size if possible
				if fontSize, err := strconv.Atoi(value); err == nil {
					scheme.Metadata.FontSize = fontSize
				}
			case "icon_enabled":
				scheme.Metadata.IconEnabled = value
			case "icon_disabled":
				scheme.Metadata.IconDisabled = value
			}
			continue
		}

		// Match mappings
		if match := mappingPattern.FindStringSubmatch(line); match != nil {
			lhs := match[1]
			rhs := strings.Fields(match[2]) // Split RHS into separate mappings
			mapping := Mapping{LHS: lhs, RHS: rhs}
			scheme.Mappings = append(scheme.Mappings, mapping)
		}
	}

	// Handle scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Write to JSON file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(scheme); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Printf("AKT file converted to JSON successfully: %s\n", outputFile)
}
