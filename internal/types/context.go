package types

import "aks.go/internal/core"

// Context maintains the state during transliteration.
// It tracks the result of the last lookup operation.
type Context struct {
	LatestLookup core.LookupResult // Tracks the result of the last lookup
}

// NewContext initializes a new Context instance with default values.
func NewContext() *Context {
	return &Context{
		LatestLookup: core.LookupResult{
			Output:   "",
			Category: "",
		},
	}
}

// Reset clears the context state, setting LatestLookup to its default values.
func (ctx *Context) Reset() {
	ctx.LatestLookup = core.LookupResult{
		Output:   "",
		Category: "",
	}
}
