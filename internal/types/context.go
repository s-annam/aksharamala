package types

import "aks.go/internal/core"

// Context maintains the state during transliteration.
type Context struct {
	LatestLookup core.LookupResult // Tracks the result of the last lookup
}

// NewContext initializes a new Context instance.
func NewContext() *Context {
	return &Context{
		LatestLookup: core.LookupResult{
			Output:   "",
			Category: "",
		},
	}
}

// Reset clears the context state.
func (ctx *Context) Reset() {
	ctx.LatestLookup = core.LookupResult{
		Output:   "",
		Category: "",
	}
}
