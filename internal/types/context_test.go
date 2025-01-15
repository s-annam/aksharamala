package types

import (
	"testing"

	"aks.go/internal/core"
)

// TestContextReset tests the Reset method of the Context struct.
// It verifies that the context state is cleared after calling Reset.
func TestContextReset(t *testing.T) {
	ctx := NewContext()

	// Update the context
	ctx.LatestLookup = core.LookupResult{
		Output:   "क",
		Category: "consonants",
	}

	// Reset the context
	ctx.Reset()

	// Validate that the context is cleared
	if ctx.LatestLookup.Output != "" || ctx.LatestLookup.Category != "" {
		t.Errorf("Expected cleared context, got %+v", ctx.LatestLookup)
	}
}
