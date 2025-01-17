package types

import "sort"

// OrderedCategories represents a map of categories that maintains a specific order
type OrderedCategories struct {
	order      []string
	categories map[string]Section
}

// OrderedCategoriesMap returns an OrderedCategories instance with categories arranged
// in a predefined order: consonants, matras, vowels, digits, vedic, others,
// followed by any remaining categories in alphabetical order.
func (s *TransliterationScheme) OrderedCategoriesMap() *OrderedCategories {
	// Define the priority order
	priorityOrder := map[string]int{
		"consonants": 0,
		"matras":     1,
		"vowels":     2,
		"digits":     3,
		"vedic":      4,
		"others":     5,
	}

	// Collect all category names
	categories := make([]string, 0, len(s.Categories))
	for category := range s.Categories {
		categories = append(categories, category)
	}

	// Sort categories based on priority order
	sort.Slice(categories, func(i, j int) bool {
		iPriority, iExists := priorityOrder[categories[i]]
		jPriority, jExists := priorityOrder[categories[j]]

		// If both categories have defined priorities, sort by priority
		if iExists && jExists {
			return iPriority < jPriority
		}
		// If only first category has priority, it comes first
		if iExists {
			return true
		}
		// If only second category has priority, it comes first
		if jExists {
			return false
		}
		// If neither has priority, sort alphabetically
		return categories[i] < categories[j]
	})

	return &OrderedCategories{
		order:      categories,
		categories: s.Categories,
	}
}

// Range iterates over the categories in order, calling fn for each key/value pair.
// The iteration stops if fn returns false.
func (oc *OrderedCategories) Range(fn func(category string, section Section) bool) {
	for _, category := range oc.order {
		if section, exists := oc.categories[category]; exists {
			if !fn(category, section) {
				break
			}
		}
	}
}

// Get returns the Section for a given category and whether it exists.
func (oc *OrderedCategories) Get(category string) (Section, bool) {
	section, exists := oc.categories[category]
	return section, exists
}

// Order returns the slice of category names in their defined order.
func (oc *OrderedCategories) Order() []string {
	return oc.order
}
