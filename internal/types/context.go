package types

import (
	"strings"
	"aks.go/internal/core"
)

// Context represents the current state of transliteration processing,
// including the latest lookup result and contextual information.
type Context struct {
	LatestLookup core.LookupResult // Tracks the result of the last lookup
	CurrentContext string // The current context marker (e.g., "M", "x")
}

// NewContext creates a new Context instance with default values.
func NewContext() *Context {
	return &Context{
		LatestLookup: core.LookupResult{
			Output:   "",
			Category: "",
		},
		CurrentContext: "",
	}
}

// Reset clears the context state, setting LatestLookup to its default values.
func (ctx *Context) Reset() {
	ctx.LatestLookup = core.LookupResult{
		Output:   "",
		Category: "",
	}
}

// ContextualRule represents a rule for modifying output based on context.
type ContextualRule struct {
	ChangePrevious bool   // [c] flag
	RequiredContext string // [M] - required context for rule to apply
	NewContext     string // [m] - context to set after applying rule
	Modification   string // The actual modification to apply
}

// ParseContextualRules parses the RHS string to extract contextual rules.
// Returns the base output and any contextual rules found.
func ParseContextualRules(rhs string) (string, []ContextualRule) {
	var rules []ContextualRule
	var baseOutput string

	// Split the RHS into parts by parentheses
	parts := strings.Split(rhs, "(")
	baseOutput = parts[0]

	for _, part := range parts[1:] {
		if !strings.HasSuffix(part, ")") {
			continue
		}
		part = strings.TrimSuffix(part, ")")

		rule := ContextualRule{}

		switch {
		case strings.HasPrefix(part, "c"):
			// Change previous character rule
			rule.ChangePrevious = true
			part = strings.TrimPrefix(part, "c")
			fallthrough

		case strings.HasPrefix(part, "M") || strings.HasPrefix(part, "x"):
			if strings.HasPrefix(part, "M") {
				// Context check rule
				rule.RequiredContext = string(part[0])
				part = part[1:]
			}
			if len(part) > 0 {
				// Context update rule
				rule.NewContext = string(part[0])
				part = part[1:]
			}
			if len(part) > 0 {
				// Remaining part is the modification
				rule.Modification = part
			}
			rules = append(rules, rule)
		}
	}

	return baseOutput, rules
}

// ApplyContextualRules applies the contextual rules to modify the output.
// Returns the modified output and any error encountered.
func (ctx *Context) ApplyContextualRules(rules []ContextualRule, builder *strings.Builder) error {
	output := builder.String()
	modified := false

	for _, rule := range rules {
		// Skip if required context doesn't match
		if rule.RequiredContext != "" && rule.RequiredContext != ctx.CurrentContext {
			continue
		}

		if rule.ChangePrevious {
			// Modify the previous character
			if len(output) > 0 {
				output = output[:len(output)-1] + rule.Modification
				modified = true
			}
		} else if rule.Modification != "" {
			output += rule.Modification
			modified = true
		}

		// Update context if specified
		if rule.NewContext != "" {
			ctx.CurrentContext = rule.NewContext
		}
	}

	if modified {
		// Clear the builder and write the modified output
		builder.Reset()
		builder.WriteString(output)
	}

	return nil
}
