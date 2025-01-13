package translit

import (
	"testing"
)

func TestContext_Reset(t *testing.T) {
	ctx := NewContext()

	// Update the context
	ctx.LatestLookup = LookupResult{
		Output:   "क",
		Category: "consonants",
	}

	// Reset the context
	ctx.Reset()

	// Validate that the context is cleared
	if ctx.LatestLookup.Output != "" || ctx.LatestLookup.Category != "" {
		t.Errorf("Context reset failed: got %+v", ctx.LatestLookup)
	} else {
		t.Log("Context reset successfully")
	}
}
