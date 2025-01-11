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

package types

import (
	"bytes"
	"encoding/json"
)

// TransliterationScheme represents a keymap for transliteration.
// Note: If you make changes to this struct, ensure that the corresponding
// fields in CompactScheme are updated accordingly.
type TransliterationScheme struct {
	Comments   []string           `json:"comments,omitempty"`
	Version    string             `json:"version"`
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	License    string             `json:"license"`
	Language   string             `json:"language"`
	Scheme     string             `json:"scheme"`
	Metadata   Metadata           `json:"metadata"`
	Categories map[string]Section `json:"categories"`
}

// Metadata contains additional configuration for a transliteration scheme.
type Metadata struct {
	Virama       string `json:"virama,omitempty"`
	FontName     string `json:"font_name,omitempty"`
	FontSize     int    `json:"font_size,omitempty"`
	IconEnabled  string `json:"icon_enabled,omitempty"`
	IconDisabled string `json:"icon_disabled,omitempty"`
}

// Section represents a category of mappings within a transliteration scheme.
type Section struct {
	Comments []string        `json:"comments,omitempty"`
	Mappings []CategoryEntry `json:"mappings"`
}

// CategoryEntry represents a mapping of LHS (input characters) to RHS (output characters).
type CategoryEntry struct {
	LHS     []string `json:"lhs"`
	RHS     []string `json:"rhs"`
	Comment string   `json:"comment,omitempty"`
}

// CompactTransliterationScheme is a temporary struct to hold the compact JSON
type CompactTransliterationScheme struct {
	Comments   []string                   `json:"comments,omitempty"`
	Version    string                     `json:"version"`
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	License    string                     `json:"license"`
	Language   string                     `json:"language"`
	Scheme     string                     `json:"scheme"`
	Metadata   Metadata                   `json:"metadata"`
	Categories map[string]json.RawMessage `json:"categories"`
}

// ToCompactTransliterationScheme converts a TransliterationScheme to a CompactTransliterationScheme
func ToCompactTransliterationScheme(scheme TransliterationScheme) (CompactTransliterationScheme, error) {
	compactCategories := make(map[string]json.RawMessage)
	for key, section := range scheme.Categories {
		var buf bytes.Buffer

		// Start array
		buf.WriteString("[")

		if len(section.Mappings) > 0 {
			buf.WriteString("\n")

			// Write each mapping
			for i, mapping := range section.Mappings {
				mappingBytes, err := json.Marshal(mapping)
				if err != nil {
					return CompactTransliterationScheme{}, err
				}

				// Add indentation and mapping
				buf.WriteString("      ")
				buf.Write(mappingBytes)

				// Add comma if not last
				if i < len(section.Mappings)-1 {
					buf.WriteString(",")
				}
				buf.WriteString("\n")
			}

			// Close array with proper indentation
			buf.WriteString("    ")
		}
		buf.WriteString("]")

		compactCategories[key] = json.RawMessage(buf.Bytes())
	}

	return CompactTransliterationScheme{
		Comments:   scheme.Comments,
		Version:    scheme.Version,
		ID:         scheme.ID,
		Name:       scheme.Name,
		License:    scheme.License,
		Language:   scheme.Language,
		Scheme:     scheme.Scheme,
		Metadata:   scheme.Metadata,
		Categories: compactCategories,
	}, nil
}
