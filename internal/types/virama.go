// Package types provides definitions and utilities for handling transliteration types,
// including virama modes used in various scripts.
package types

import (
	"fmt"
	"strconv"
	"strings"
)

// ViramaMode represents different modes of virama handling in transliteration.
type ViramaMode int

// UnknownMode indicates an unrecognized virama mode.
// NormalMode indicates standard virama handling.
// SmartMode indicates intelligent virama handling based on context.
const (
	UnknownMode ViramaMode = iota
	NormalMode
	SmartMode
)

// viramaModeStrings maps string representations of virama modes to their corresponding ViramaMode values.
var viramaModeStrings = map[string]ViramaMode{
	"normal": NormalMode,
	"smart":  SmartMode,
}

// ParseViramaMode converts a string representation of a virama mode to its corresponding ViramaMode value.
// If the input string does not match any known modes, it returns UnknownMode.
func ParseViramaMode(mode string) ViramaMode {
	if v, ok := viramaModeStrings[mode]; ok {
		return v
	}
	return UnknownMode
}

// String returns a string representation of the ViramaMode value.
func (v ViramaMode) String() string {
	for k, value := range viramaModeStrings {
		if value == v {
			return k
		}
	}
	return "unknown"
}

// splitAndTrim splits a string by commas and trims whitespace from each part.
// It returns a slice of trimmed strings.
func splitAndTrim(s string) []string {
	// Split the string by commas and trim white spaces
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ','
	})

	// Trim white spaces from each part
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}

// ParseVirama parses the virama metadata string and returns the corresponding virama character and mode.
// It returns an error if the input is invalid.
func ParseVirama(metadata string) (string, ViramaMode, error) {
	parts := splitAndTrim(metadata)
	if len(parts) != 2 {
		return "", UnknownMode, fmt.Errorf("invalid or missing virama metadata: %s", metadata)
	}

	viramaString := parts[0]
	mode := ParseViramaMode(parts[1])
	if mode == UnknownMode {
		return "", UnknownMode, fmt.Errorf("invalid virama mode: %s", parts[1])
	}

	// Handle hexadecimal Unicode value, if present
	if strings.HasPrefix(viramaString, "0x") {
		codePoint, err := strconv.ParseUint(viramaString[2:], 16, 32)
		if err != nil {
			return "", UnknownMode, fmt.Errorf("invalid Unicode code point: %v", err)
		}
		return string(rune(codePoint)), mode, nil
	} else if viramaString == "" {
		return "", UnknownMode, fmt.Errorf("empty virama: %s", metadata)
	}

	return viramaString, mode, nil
}
