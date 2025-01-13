package translit

// Context maintains the state during transliteration.
type Context struct {
	LatestLookup LookupResult // Tracks the result of the last lookup
}

// NewContext initializes a new Context instance.
func NewContext() *Context {
	return &Context{
		LatestLookup: LookupResult{
			Output:   "",
			Category: "",
		},
	}
}

// Reset clears the context state.
func (ctx *Context) Reset() {
	ctx.LatestLookup = LookupResult{
		Output:   "",
		Category: "",
	}
}
