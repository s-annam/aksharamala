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
	"encoding/json"
	"fmt"
	"strings"
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

func normalizeComment(comment string) string {
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "=*=")
	comment = strings.TrimSuffix(comment, "=*=")
	comment = strings.TrimSpace(comment)
	return comment
}

// ToCompactTransliterationScheme converts a TransliterationScheme to a CompactTransliterationScheme
func ToCompactTransliterationScheme(scheme TransliterationScheme) (CompactTransliterationScheme, error) {
	compactCategories := make(map[string]json.RawMessage)

	for category, section := range scheme.Categories {
		// Normalize comments in mappings
		for i := range section.Mappings {
			section.Mappings[i].Comment = normalizeComment(section.Mappings[i].Comment)
		}

		// Convert section to JSON
		sectionJSON, err := json.Marshal(section.Mappings)
		if err != nil {
			return CompactTransliterationScheme{}, err
		}
		compactCategories[category] = sectionJSON
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

// UnmarshalJSON customizes JSON unmarshaling for TransliterationScheme
func (s *TransliterationScheme) UnmarshalJSON(data []byte) error {
	var compact CompactTransliterationScheme
	if err := json.Unmarshal(data, &compact); err != nil {
		return err
	}

	// Copy the simple fields
	s.Comments = compact.Comments
	s.Version = compact.Version
	s.ID = compact.ID
	s.Name = compact.Name
	s.License = compact.License
	s.Language = compact.Language
	s.Scheme = compact.Scheme
	s.Metadata = compact.Metadata

	// Initialize the Categories map
	s.Categories = make(map[string]Section)

	// Process each category
	for name, rawEntries := range compact.Categories {
		var mappings []CategoryEntry
		if err := json.Unmarshal(rawEntries, &mappings); err != nil {
			return err
		}
		s.Categories[name] = Section{
			Mappings: mappings,
		}
	}

	return nil
}

// Validate checks the integrity of the transliteration scheme.
func (s *TransliterationScheme) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("keymap is missing 'id'")
	}
	if len(s.Categories) == 0 {
		return fmt.Errorf("keymap '%s' has no categories", s.ID)
	}
	for category, section := range s.Categories {
		for _, mapping := range section.Mappings {
			if len(mapping.LHS) == 0 {
				return fmt.Errorf("category '%s' in keymap '%s' has an empty 'LHS'", category, s.ID)
			}
			if len(mapping.RHS) == 0 {
				return fmt.Errorf("category '%s' in keymap '%s' has an empty 'RHS'", category, s.ID)
			}
		}
	}
	return nil
}
