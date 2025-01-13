package types

import "fmt"

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

// Lookup searches for a match in the mappings by LHS.
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

// Validate checks the integrity of all mappings.
func (m *Mappings) Validate() error {
	for _, entry := range m.entries {
		if len(entry.LHS) == 0 {
			return fmt.Errorf("Mapping with empty LHS found")
		}
		if len(entry.RHS) == 0 {
			return fmt.Errorf("Mapping with empty RHS found")
		}
	}
	return nil
}

// All returns all mappings in the collection.
func (m *Mappings) All() []Mapping {
	return m.entries
}

// Entries returns the underlying slice of Mapping.
func (m *Mappings) Entries() []Mapping {
	return m.entries
}

// Get returns the Mapping at the specified index.
func (m *Mappings) Get(index int) (Mapping, error) {
	if index < 0 || index >= len(m.entries) {
		return Mapping{}, fmt.Errorf("index out of bounds")
	}
	return m.entries[index], nil
}

// NewMappings initializes a Mappings instance from a slice of Mapping.
func NewMappings(entries []Mapping) Mappings {
	return Mappings{entries: entries}
}
