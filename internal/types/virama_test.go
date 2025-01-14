package types

import (
	"testing"
)

// TestParseVirama tests the ParseVirama method of the Aksharamala struct.
// It verifies that the virama parsing works correctly for various input cases.
func TestParseVirama(t *testing.T) {
	tests := []struct {
		metadata       string
		expectedVirama string
		expectedMode   ViramaMode
		shouldError    bool
	}{
		{"0x094D, smart", "्", SmartMode, false},
		{"्, normal", "्", NormalMode, false},
		{"abcd, smart", "abcd", SmartMode, false},
		{"0xZZZZ, smart", "", UnknownMode, true},
		{"no_comma", "", UnknownMode, true},
	}

	for _, test := range tests {
		r, mode, err := ParseVirama(test.metadata)
		if test.shouldError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.metadata)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.metadata, err)
			continue
		}

		if r != test.expectedVirama || mode != test.expectedMode {
			t.Errorf("For input '%s', expected '%s' and mode '%s', but got '%s' and mode '%s'",
				test.metadata, test.expectedVirama, test.expectedMode, r, mode)
		}
	}
}
