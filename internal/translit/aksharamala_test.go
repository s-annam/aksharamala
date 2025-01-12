package translit

import (
	"testing"

	"aks.go/internal/types"
)

func TestNewAksharamala(t *testing.T) {
	tests := []struct {
		name    string
		scheme  *types.TransliterationScheme
		wantErr bool
	}{
		{
			name:    "nil scheme",
			scheme:  nil,
			wantErr: true,
		},
		{
			name: "valid scheme",
			scheme: &types.TransliterationScheme{
				Categories: map[string]types.Section{
					"consonants": {
						Mappings: []types.CategoryEntry{
							{LHS: []string{"k"}, RHS: []string{"क"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid mapping",
			scheme: &types.TransliterationScheme{
				Categories: map[string]types.Section{
					"consonants": {
						Mappings: []types.CategoryEntry{
							{LHS: []string{}, RHS: []string{}}, // empty mapping
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAksharamala(tt.scheme)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAksharamala() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewAksharamala() returned nil but wanted valid instance")
			}
		})
	}
}

func TestAksharamala_Transliterate(t *testing.T) {
	// Basic test scheme with some Hindi mappings
	scheme := &types.TransliterationScheme{
		Categories: map[string]types.Section{
			"consonants": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"k"}, RHS: []string{"क"}},
					{LHS: []string{"kh"}, RHS: []string{"ख"}},
					{LHS: []string{"g"}, RHS: []string{"ग"}},
				},
			},
			"vowels": {
				Mappings: []types.CategoryEntry{
					{LHS: []string{"a"}, RHS: []string{"अ"}},
					{LHS: []string{"aa"}, RHS: []string{"आ"}},
				},
			},
		},
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "single consonant",
			input: "k",
			want:  "क",
		},
		{
			name:  "consonant combination",
			input: "kh",
			want:  "ख",
		},
		{
			name:  "multiple characters",
			input: "kag",
			want:  "कग",
		},
		{
			name:  "english mode toggle",
			input: "k#hello#k",
			want:  "कhelloक",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "unmapped character",
			input: "x",
			want:  "x",
		},
		{
			name:  "mixed mapped and unmapped",
			input: "k123g",
			want:  "क123ग",
		},
	}

	aks, err := NewAksharamala(scheme)
	if err != nil {
		t.Fatalf("Failed to create Aksharamala: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aks.Transliterate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transliterate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Transliterate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer(t *testing.T) {
	b := NewBuffer()

	// Test initial state
	if b.Len() != 0 {
		t.Errorf("New buffer length = %d, want 0", b.Len())
	}

	// Test append
	b.Append('a')
	b.Append('b')
	if b.Len() != 2 {
		t.Errorf("Buffer length after append = %d, want 2", b.Len())
	}

	// Test string conversion
	if s := b.String(); s != "ab" {
		t.Errorf("Buffer.String() = %s, want ab", s)
	}

	// Test first character
	if f := b.First(); f != 'a' {
		t.Errorf("Buffer.First() = %c, want a", f)
	}

	// Test remove first
	b.RemoveFirst()
	if s := b.String(); s != "b" {
		t.Errorf("After RemoveFirst(), String() = %s, want b", s)
	}

	// Test clear
	b.Clear()
	if b.Len() != 0 {
		t.Errorf("After Clear(), length = %d, want 0", b.Len())
	}
}

func TestContext(t *testing.T) {
	ctx := NewContext()

	// Test initial state
	if ctx.InEnglishMode {
		t.Error("New context should not be in English mode")
	}

	// Test English mode toggle
	ctx.ToggleEnglishMode()
	if !ctx.InEnglishMode {
		t.Error("Context should be in English mode after toggle")
	}
	ctx.ToggleEnglishMode()
	if ctx.InEnglishMode {
		t.Error("Context should not be in English mode after second toggle")
	}

	// Test output update
	ctx.UpdateWithOutput("क")
	if ctx.LastChar != 'क' {
		t.Errorf("LastChar = %c, want क", ctx.LastChar)
	}
}
