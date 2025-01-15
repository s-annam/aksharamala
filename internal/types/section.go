package types

import (
	"aks.go/internal/core"
)

// Section represents a category of mappings within a transliteration scheme.
// It contains comments and a collection of mappings that define the relationships
// between left-hand side (LHS) and right-hand side (RHS) elements.
type Section struct {
	Comments []string      `json:"comments,omitempty"` // Optional comments about the section
	Mappings core.Mappings `json:"mappings"`           // Mappings for the section
}

// NewSection initializes and returns a new Section with empty comments and mappings.
func NewSection() Section {
	return Section{
		Comments: []string{},
		Mappings: core.NewMappings([]core.Mapping{}),
	}
}

// GetOrCreate returns an existing section for the given category.
// If the section does not exist, it creates a new one.
func GetOrCreate(scheme TransliterationScheme, categoryName string) *Section {
	if section, exists := scheme.Categories[categoryName]; exists {
		return &section
	}
	section := NewSection()
	scheme.Categories[categoryName] = section
	return &section
}

// AddMapping adds a full mapping (LHS -> RHS) to the section.
// It takes LHS and RHS as slices of strings and an optional comment.
func (s *Section) AddMapping(lhs []string, rhs []string, comment string) {
	s.Mappings.Add(lhs, rhs, comment)
}

// AppendLHSToMapping appends an LHS value to the last mapping in the section.
func (s *Section) AppendLHSToMapping(lastMapping *core.Mapping, lhs string) {
	if lastMapping == nil {
		return
	}

	entries := s.Mappings.All()
	for i, entry := range entries {
		if entry.LHS[0] == lastMapping.LHS[0] {
			entries[i].LHS = append(entries[i].LHS, lhs)
			s.Mappings = core.NewMappings(entries)
			break
		}
	}
}

// GetMappings returns all mappings in the section.
// It returns a slice of core.Mapping.
func (s *Section) GetMappings() []core.Mapping {
	return s.Mappings.All()
}
