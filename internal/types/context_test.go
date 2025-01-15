package types

import (
	"strings"
	"testing"

	"aks.go/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestContextReset(t *testing.T) {
	ctx := NewContext()
	ctx.LatestLookup = core.LookupResult{
		Output:   "test",
		Category: "test",
	}
	ctx.CurrentContext = "M"

	ctx.Reset()

	if ctx.CurrentContext != "" {
		t.Errorf("Expected empty context, got %s", ctx.CurrentContext)
	}
	if ctx.LatestLookup.Output != "" || ctx.LatestLookup.Category != "" {
		t.Errorf("Expected cleared context, got %+v", ctx.LatestLookup)
	}
}

func TestIsNextWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		position int
		expected bool
	}{
		{
			name:     "End of string",
			input:    "test",
			position: 3,
			expected: true,
		},
		{
			name:     "Next char is space",
			input:    "test test",
			position: 3,
			expected: true,
		},
		{
			name:     "Middle of word",
			input:    "testing",
			position: 2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			ctx.Input = tt.input
			ctx.Position = tt.position
			result := ctx.IsSeparator()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseContextualRulesWithWhitespace(t *testing.T) {
	// In TeluguRts.aksj, RHS is commonly an array of strings where:
	// 1. First element: Base character with context rule, e.g., "మ(M)" - outputs మ and sets context to M
	// 2. Second element: Contextual rule for sunna, e.g., "(W)ం" - if whitespace follows, outputs anusvara (sunna)
	//
	// Examples from TeluguRts.aksj:
	// {"lhs":["m"],"rhs":["మ(M)","(W)ం"]} - 'm' maps to మ with context M, followed by anusvara if whitespace
	// {"lhs":["n"],"rhs":["న(M)","(W)ం"]} - 'n' maps to న with context M, followed by anusvara if whitespace
	tests := []struct {
		name          string
		rhs           string
		expectedBase  string
		expectedRules []ContextualRule
	}{
		{
			name:         "Base character with context (first RHS element)",
			rhs:          "మ(M)",
			expectedBase: "మ",
			expectedRules: []ContextualRule{
				{
					RequiredContext: "M",
				},
			},
		},
		{
			name:         "Sunna with whitespace requirement (second RHS element)",
			rhs:          "(W)ం",
			expectedBase: "",
			expectedRules: []ContextualRule{
				{
					WhitespaceRequired: true,
					Modification:       "ం",
				},
			},
		},
		{
			name:         "Consonant with context and modification",
			rhs:          "క(c)(M)(x)ం",
			expectedBase: "క",
			expectedRules: []ContextualRule{
				{
					ChangePrevious: true,
				},
				{
					RequiredContext: "M",
				},
				{
					NewContext:   "x",
					Modification: "ం",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, rules := ParseContextualRules(tt.rhs)
			assert.Equal(t, tt.expectedBase, base)
			assert.Equal(t, len(tt.expectedRules), len(rules))

			for i, expectedRule := range tt.expectedRules {
				if i < len(rules) {
					assert.Equal(t, expectedRule.WhitespaceRequired, rules[i].WhitespaceRequired, "WhitespaceRequired mismatch at index %d", i)
					assert.Equal(t, expectedRule.Modification, rules[i].Modification, "Modification mismatch at index %d", i)
					assert.Equal(t, expectedRule.RequiredContext, rules[i].RequiredContext, "RequiredContext mismatch at index %d", i)
					assert.Equal(t, expectedRule.ChangePrevious, rules[i].ChangePrevious, "ChangePrevious mismatch at index %d", i)
					assert.Equal(t, expectedRule.NewContext, rules[i].NewContext, "NewContext mismatch at index %d", i)
				}
			}
		})
	}
}

func TestApplyContextualRulesWithWhitespace(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		position       int
		rules          []ContextualRule
		initialOutput  string
		expectedOutput string
	}{
		{
			name:     "Apply whitespace rule at end",
			input:    "test",
			position: 3,
			rules: []ContextualRule{
				{
					WhitespaceRequired: true,
					Modification:       "ం",
				},
			},
			initialOutput:  "మ",
			expectedOutput: "మం",
		},
		{
			name:     "Skip whitespace rule in middle",
			input:    "testing",
			position: 2,
			rules: []ContextualRule{
				{
					WhitespaceRequired: true,
					Modification:       "ం",
				},
			},
			initialOutput:  "మ",
			expectedOutput: "మ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			ctx.Input = tt.input
			ctx.Position = tt.position

			var builder strings.Builder
			builder.WriteString(tt.initialOutput)

			err := ctx.ApplyContextualRules(tt.rules, &builder)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, builder.String())
		})
	}
}
