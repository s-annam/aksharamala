package main

import (
	"encoding/json"
	"sort"
	"strings"

	"aks.go/internal/types"
)

// FormatSchemeJSON converts a CompactTransliterationScheme to a formatted JSON string.
// It takes a CompactTransliterationScheme and returns the JSON representation as a string.
func FormatSchemeJSON(scheme types.CompactTransliterationScheme) (string, error) {
	// FormatSchemeJSON takes a CompactTransliterationScheme and returns its JSON representation as a string.
	var output strings.Builder
	output.WriteString("{\n")

	// Write metadata fields
	if err := writeMetadataFields(&output, scheme); err != nil {
		return "", err
	}

	// Write categories
	if err := writeCategories(&output, scheme.Categories); err != nil {
		return "", err
	}

	output.WriteString("}\n")
	return output.String(), nil
}

// writeMetadataFields writes the metadata fields of the scheme to the provided string builder.
// It takes a pointer to strings.Builder and a CompactTransliterationScheme.
// The function writes the metadata fields in the following order: comments, version, id, name, license, language, scheme, and metadata.
func writeMetadataFields(w *strings.Builder, scheme types.CompactTransliterationScheme) error {
	// Comments with default formatting (one per line)
	w.WriteString(`  "comments": [`)
	for i, comment := range scheme.Comments {
		w.WriteString("\n    ")
		commentJSON, err := json.Marshal(comment)
		if err != nil {
			return err
		}
		w.Write(commentJSON)
		if i < len(scheme.Comments)-1 {
			w.WriteString(",")
		}
	}
	w.WriteString("\n  ],\n")

	// Write other fields
	writeStringField(w, "version", scheme.Version)
	writeStringField(w, "id", scheme.ID)
	writeStringField(w, "name", scheme.Name)
	writeStringField(w, "license", scheme.License)
	writeStringField(w, "language", scheme.Language)
	writeStringField(w, "scheme", scheme.Scheme)

	// Metadata object
	w.WriteString(`  "metadata": `)
	metadataJSON, err := json.Marshal(scheme.Metadata)
	if err != nil {
		return err
	}
	w.Write(metadataJSON)
	w.WriteString(",\n")

	return nil
}

// writeStringField writes a string field to the provided string builder.
// It takes a pointer to strings.Builder, a field name, and a field value.
// The function writes the field in the format "field": "value".
func writeStringField(w *strings.Builder, field, value string) {
	w.WriteString(`  "`)
	w.WriteString(field)
	w.WriteString(`": "`)
	w.WriteString(value)
	w.WriteString(`",`)
	w.WriteString("\n")
}

// writeCategories writes the categories of the scheme to the provided string builder.
// It takes a pointer to strings.Builder and a map of categories.
// The function writes the categories in sorted order, with each category containing a list of mappings.
func writeCategories(w *strings.Builder, categories map[string]json.RawMessage) error {
	w.WriteString(`  "categories": {`)

	// Get sorted category names
	categoryNames := make([]string, 0, len(categories))
	for name := range categories {
		categoryNames = append(categoryNames, name)
	}
	sort.Strings(categoryNames)

	// Write categories in sorted order
	for i, category := range categoryNames {
		mappings := categories[category]
		if i > 0 {
			w.WriteString(",")
		}
		w.WriteString("\n    \"")
		w.WriteString(category)
		w.WriteString("\": [\n")

		var entries []map[string]interface{}
		if err := json.Unmarshal(mappings, &entries); err != nil {
			return err
		}

		// Write each mapping
		for i, entry := range entries {
			w.WriteString("      {")

			// lhs first
			if lhs, ok := entry["lhs"]; ok {
				lhsJSON, err := json.Marshal(lhs)
				if err != nil {
					return err
				}
				w.WriteString(`"lhs":`)
				w.Write(lhsJSON)
			}

			// rhs second
			if rhs, ok := entry["rhs"]; ok {
				w.WriteString(`,"rhs":`)
				rhsJSON, err := json.Marshal(rhs)
				if err != nil {
					return err
				}
				w.Write(rhsJSON)
			}

			// comment last, if it exists
			if comment, ok := entry["comment"]; ok && comment != nil {
				w.WriteString(`,"comment":`)
				commentJSON, err := json.Marshal(comment)
				if err != nil {
					return err
				}
				w.Write(commentJSON)
			}

			w.WriteString("}")
			if i < len(entries)-1 {
				w.WriteString(",")
			}
			w.WriteString("\n")
		}
		w.WriteString("    ]")
	}
	w.WriteString("\n  }\n")
	return nil
}
