package types

import (
	"strings"
	"unicode"

	"aks.go/internal/core"
)

// Context represents the current state of transliteration processing,
// including the latest lookup result and contextual information.
type Context struct {
	LatestLookup   core.LookupResult // Tracks the result of the last lookup
	CurrentContext string            // The current context marker (e.g., "M", "x")
	Input          string            // The full input string being processed
	Position       int               // Current position in the input string
}

// NewContext creates a new Context instance with default values.
func NewContext() *Context {
	return &Context{
		LatestLookup:   core.LookupResult{},
		CurrentContext: "",
		Input:          "",
		Position:       0,
	}
}

// Reset clears the context state, setting all fields to their default values.
func (ctx *Context) Reset() {
	ctx.LatestLookup = core.LookupResult{}
	ctx.CurrentContext = ""
	ctx.Input = ""
	ctx.Position = 0
}

// ContextualRule represents a rule for modifying output based on context.
type ContextualRule struct {
	ChangePrevious     bool   // (c) flag
	RequiredContext    string // (M) - required context for rule to apply
	NewContext         string // (m) - context to set after applying rule
	WhitespaceRequired bool   // (W) flag - requires next char to be whitespace or EOS
	Modification       string // The actual modification to apply
}

// IsSeparator checks if the next character in the input is a separator (whitespace, punctuation, etc.)
// or if we're at the end of input. A separator is any character that's not part of a consonant or vowel.
func (ctx *Context) IsSeparator() bool {
	// Get all runes from the input
	runes := []rune(ctx.Input)

	// Check if we're at the last position
	if ctx.Position+1 >= len(runes) {
		return true
	}

	// Get the next rune
	nextRune := runes[ctx.Position+1]

	// Check if it's a separator (whitespace, punctuation, or non-letter)
	// IsLetter() alone is not enough as it doesn't include combining marks
	return unicode.IsSpace(nextRune) || unicode.IsPunct(nextRune) ||
		!(unicode.IsLetter(nextRune) || unicode.IsMark(nextRune))
}

// ShouldApplyRule determines if a contextual rule should be applied based on all conditions
func (ctx *Context) ShouldApplyRule(rule ContextualRule) bool {
	// Check context requirement
	if rule.RequiredContext != "" && rule.RequiredContext != ctx.CurrentContext {
		return false
	}

	// Check whitespace requirement
	if rule.WhitespaceRequired && !ctx.IsSeparator() {
		return false
	}

	return true
}

// parseRule extracts a single contextual rule from a rule string
func parseRule(ruleStr string) (ContextualRule, bool) {
	var rule ContextualRule
	if len(ruleStr) < 1 {
		return rule, false
	}

	switch ruleStr[0] {
	case 'W':
		rule.WhitespaceRequired = true
		rule.Modification = ruleStr[1:]
	case 'c':
		rule.ChangePrevious = true
		rule.Modification = ruleStr[1:]
	case 'M':
		rule.RequiredContext = ruleStr
	case 'x':
		rule.NewContext = ruleStr
	default:
		return rule, false
	}
	return rule, true
}

// ParseContextualRules parses the RHS string to extract contextual rules.
// Returns the base output and any contextual rules found.
func ParseContextualRules(rhs string) (string, []ContextualRule) {
	var rules []ContextualRule

	// Find the first rule marker
	firstRule := strings.Index(rhs, "(")
	if firstRule == -1 {
		return rhs, nil // No rules found
	}

	// Extract base output (text before first rule)
	baseOutput := rhs[:firstRule]

	// Split remaining string into segments
	segments := strings.Split(rhs[firstRule:], ")")

	for i, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		// Extract rule content from parentheses
		ruleStart := strings.Index(segment, "(")
		if ruleStart == -1 {
			// This is a modification for the previous rule
			if len(rules) > 0 {
				rules[len(rules)-1].Modification = segment
			}
			continue
		}

		ruleContent := segment[ruleStart+1:]
		if rule, valid := parseRule(ruleContent); valid {
			// Check for modification text before next rule
			if i < len(segments)-1 {
				nextRuleStart := strings.Index(segments[i+1], "(")
				if nextRuleStart > 0 {
					rule.Modification = segments[i+1][:nextRuleStart]
				}
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
		// Check if rule should be applied
		if !ctx.ShouldApplyRule(rule) {
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
