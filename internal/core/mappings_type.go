package core

import (
	"fmt"
)

// Mappings represents a collection of transliteration mappings.
type Mappings struct {
	entries []Mapping
}

// Mapping represents a single transliteration mapping.
type Mapping struct {
	LHS     []string `json:"lhs"`
	RHS     []string `json:"rhs"`
	Comment string   `json:"comment,omitempty"`
}

// Add adds a new mapping to the collection.
func (m *Mappings) Add(lhs []string, rhs []string, comment string) {
	m.entries = append(m.entries, Mapping{LHS: lhs, RHS: rhs, Comment: comment})
}

// FindByLHS searches for the first mapping that matches the given LHS.
func (m *Mappings) FindByLHS(lhs string) (Mapping, bool) {
	for _, entry := range m.entries {
		for _, lhsCandidate := range entry.LHS {
			if lhsCandidate == lhs {
				return entry, true
			}
		}
	}
	return Mapping{}, false
}

// NormalizeComments applies a normalization function to all comments.
func (m *Mappings) NormalizeComments(normalizer func(string) string) {
	for i := range m.entries {
		m.entries[i].Comment = normalizer(m.entries[i].Comment)
	}
}

// Validate checks the integrity of all mappings.
func (m *Mappings) Validate(category, schemeID string) error {
	for _, entry := range m.entries {
		if len(entry.LHS) == 0 {
			return fmt.Errorf("category '%s' in keymap '%s' has an empty 'LHS'", category, schemeID)
		}
		if len(entry.RHS) == 0 {
			return fmt.Errorf("category '%s' in keymap '%s' has an empty 'RHS'", category, schemeID)
		}
	}
	return nil
}

// ToLookupTable adds mappings to a LookupTable with the specified category.
func (m *Mappings) ToLookupTable(category string, table LookupTable) {
	for _, entry := range m.entries {
		for _, lhs := range entry.LHS {
			table[lhs] = LookupResult{Output: entry.RHS[0], Category: category}
		}
	}
}

// Entries returns the underlying slice of Mapping.
func (m *Mappings) Entries() []Mapping {
	return m.entries
}

// NewMappings initializes a Mappings instance from a slice of Mapping.
func NewMappings(entries []Mapping) Mappings {
	return Mappings{entries: entries}
}

func (m *Mappings) Lookup(lhs string) (string, bool) {
	for _, entry := range m.entries {
		for _, lhsCandidate := range entry.LHS {
			if lhsCandidate == lhs {
				return entry.RHS[0], true
			}
		}
	}
	return "", false
}

func (m *Mappings) All() []Mapping {
	return m.entries
}

func (m *Mappings) ValidateAll(category, schemeID string) error {
	for _, entry := range m.entries {
		if len(entry.LHS) == 0 {
			return fmt.Errorf("category '%s' in keymap '%s' has an empty 'LHS'", category, schemeID)
		}
		if len(entry.RHS) == 0 {
			return fmt.Errorf("category '%s' in keymap '%s' has an empty 'RHS'", category, schemeID)
		}
	}
	return nil
}
