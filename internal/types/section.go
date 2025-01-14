package types

import (
	"fmt"

	"aks.go/internal/core"
)

// Section represents a category of mappings within a transliteration scheme.
type Section struct {
	Comments []string      `json:"comments,omitempty"`
	Mappings core.Mappings `json:"mappings"`
}

// NewSection initializes and returns a new Section.
func NewSection() Section {
	return Section{
		Comments: []string{},
		Mappings: core.NewMappings([]core.Mapping{}),
	}
}

// GetOrCreate returns an existing section for the given category.
func GetOrCreate(scheme TransliterationScheme, categoryName string) *Section {
	if section, exists := scheme.Categories[categoryName]; exists {
		return &section
	}
	section := NewSection()
	scheme.Categories[categoryName] = section
	return &section
}

// AddMapping adds a full mapping (LHS -> RHS) to the section.
func (s *Section) AddMapping(lhs []string, rhs []string, comment string) {
	s.Mappings.Add(lhs, rhs, comment)
}

// AppendLHSToMapping appends an LHS-only pattern to the last mapping in the section.
// If no last mapping exists, it creates a new mapping with the LHS-only pattern.
func (s *Section) AppendLHSToMapping(lastMapping *core.Mapping, lhs string) error {
	entries := s.Mappings.All()
	if len(entries) == 0 {
		return fmt.Errorf("no previous mapping found to update")
	}

	// Append LHS to the last mapping
	lastMapping.LHS = append(lastMapping.LHS, lhs)

	// Find last mapping in the section
	for i, entry := range s.Mappings.Entries() {
		if entry.LHS[0] == lastMapping.LHS[0] {
			s.Mappings.Entries()[i] = *lastMapping
			return nil
		}
	}

	return nil
}

// GetMappings returns all mappings in the section.
func (s *Section) GetMappings() []core.Mapping {
	return s.Mappings.All()
}