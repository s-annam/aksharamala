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

	"aks.go/internal/core"
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
	Comments []string      `json:"comments,omitempty"`
	Mappings core.Mappings `json:"mappings"`
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
		var mappings []core.Mapping
		if err := json.Unmarshal(rawEntries, &mappings); err != nil {
			return err
		}
		s.Categories[name] = Section{
			Mappings: core.NewMappings(mappings),
		}
	}

	return nil
}

// Validate checks the integrity of the transliteration scheme.
func (s *TransliterationScheme) Validate() error {
	missingFields := []string{}

	// Validate mandatory fields
	if s.ID == "" {
		missingFields = append(missingFields, "id")
		s.ID = "unknown_id" // Assign default if required
	}
	if s.Name == "" {
		missingFields = append(missingFields, "name")
		s.Name = "Unnamed Transliteration"
	}
	if s.Language == "" {
		missingFields = append(missingFields, "language")
		s.Language = "unknown_language"
	}
	if s.Scheme == "" {
		missingFields = append(missingFields, "scheme")
		s.Scheme = "unknown_scheme"
	}

	// Check categories and mappings
	if len(s.Categories) == 0 {
		return fmt.Errorf("keymap '%s' has no categories", s.ID)
	}
	for category, section := range s.Categories {
		if err := section.Mappings.ValidateAll(category, s.ID); err != nil {
			return err
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("mandatory fields missing: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

// BuildLookupTable constructs a precomputed lookup table from a transliteration scheme.
func BuildLookupTable(scheme *TransliterationScheme) core.LookupTable {
	table := make(core.LookupTable)
	for category, section := range scheme.Categories {
		for _, mapping := range section.Mappings.All() {
			for _, lhs := range mapping.LHS {
				table[lhs] = core.LookupResult{
					Output:   mapping.RHS[0],
					Category: category,
				}
			}
		}
	}
	return table
}

// ToCompactTransliterationScheme converts a TransliterationScheme to a CompactTransliterationScheme
func ToCompactTransliterationScheme(scheme TransliterationScheme) (CompactTransliterationScheme, error) {
	compactCategories := make(map[string]json.RawMessage)
	var errList []error

	scheme.IterateCategories(func(category string, section Section) {
		section.Mappings.NormalizeComments(core.NormalizeComment)
		sectionJSON, err := json.Marshal(section.Mappings.Entries())
		if err != nil {
			errList = append(errList, fmt.Errorf("failed to marshal category '%s': %w", category, err))
			return
		}
		compactCategories[category] = sectionJSON
	})

	if len(errList) > 0 {
		return CompactTransliterationScheme{}, fmt.Errorf("errors encountered: %v", errList)
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

// IterateCategories performs an action on each category and section.
func (s *TransliterationScheme) IterateCategories(action func(string, Section)) {
	for category, section := range s.Categories {
		action(category, section)
	}
}
