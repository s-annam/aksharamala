package translit

import (
	"fmt"
	"testing"

	"aks.go/internal/keymap"
)

func TestReversliterate(t *testing.T) {
	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps("../../keymaps"); err != nil {
		fmt.Printf("Failed to load keymaps: %v\n", err)
		return
	}

	aks := NewAksharamala(store)

	// Test smart mode (RHindi scheme)
	smartTests := []struct {
		input    string
		expected string
		desc     string
	}{
		// Basic consonants (should not include 'a' in smart mode)
		{"क", "k", "single consonant"},
		{"ख", "kh", "single aspirated consonant"},
		{"ग", "g", "single voiced consonant"},

		// Conjuncts and clusters
		{"क्ष", "x", "special conjunct ksha"},
		{"ज्ञ", "GY", "special conjunct gya"},
		{"क्क", "kk", "geminated consonant"},
		{"त्त", "tt", "geminated dental"},

		// Words with implicit 'a'
		{"नमस्ते", "namaste", "basic word with conjunct"},
		{"प्रेम", "prem", "word with ra-conjunct"},
		{"क्षमा", "xamaa", "word starting with special conjunct"},
		{"ज्ञान", "GYaan", "word with special conjunct and long vowel"},

		// Complex cases
		{"संस्कृत", "saMskRRit", "complex word with anusvara and vocalic r"},
		{"हिंदी", "hiMdI", "word with anusvara and long vowel"},
	}

	// Test RHindi (smart mode)
	if err := aks.SetActiveKeymap("rhindi"); err != nil {
		t.Fatalf("Failed to set active keymap: %v", err)
	}

	for _, test := range smartTests {
		t.Run(fmt.Sprintf("Smart_%s(%s)", test.input, test.desc), func(t *testing.T) {
			output, err := aks.Reversliterate(test.input)
			if err != nil {
				t.Errorf("Error reversliterating '%s': %v", test.input, err)
				return
			}
			if output != test.expected {
				t.Errorf("For input '%s' (%s): expected '%s', got '%s'",
					test.input, test.desc, test.expected, output)
			} else {
				t.Logf("Successfully reversliterated '%s' to '%s' (%s)",
					test.input, output, test.desc)
			}
		})
	}

	// Test normal mode (RSanskrit scheme)
	normalTests := []struct {
		input    string
		expected string
		desc     string
	}{
		// Basic consonants (should include 'a' in normal mode)
		{"क", "ka", "single consonant"},
		{"ख", "kha", "single aspirated consonant"},
		{"ग", "ga", "single voiced consonant"},

		// Special consonants and clusters
		{"क्ष", "xa", "special conjunct ksha"},
		{"ज्ञ", "GYa", "special conjunct gya"},
		{"क्क", "kka", "geminated initial k"}, // Different from smart mode!
		{"त्त", "tta", "geminated initial t"}, // Different from smart mode!

		// Sanskrit words
		{"धर्म", "dharma", "basic word with ra-conjunct"},
		{"कर्म", "karma", "word with ra-conjunct"},
		{"सत्य", "satya", "word with ya-conjunct"},
		{"अद्य", "adya", "word starting with vowel"},

		// Words with special vowels
		{"कृष्ण", "kRRiShNa", "word with vocalic r"},
		{"दृष्टि", "dRRiShTi", "word with vocalic r and conjunct"},

		// Complex cases
		{"संस्कृत", "saMskRRita", "word with anusvara and vocalic r"}, // Note: no final 'a'
		{"सन्धि", "sandhi", "word with homorganic nasal"},

		// Special Sanskrit features
		{"अग्निः", "agniH", "word with visarga"},
		{"देवाः", "devaaH", "word with long A and visarga"},
		{"गङ्गा", "ga~Ngaa", "word with velar nasal"},

		// Vedic Sanskrit specific (if supported)
		{"वॢ", "vLLi", "vocalic l"},
		{"वॣ", "vLLI", "long vocalic l"},
		{"ॐ", "_AUM_", "sacred syllable OM"},
	}

	// Test RSanskrit (normal mode)
	if err := aks.SetActiveKeymap("rsanskrit"); err != nil {
		t.Fatalf("Failed to set active keymap: %v", err)
	}

	for _, test := range normalTests {
		t.Run(fmt.Sprintf("Normal_%s(%s)", test.input, test.desc), func(t *testing.T) {
			output, err := aks.Reversliterate(test.input)
			if err != nil {
				t.Errorf("Error reversliterating '%s': %v", test.input, err)
				return
			}
			if output != test.expected {
				t.Errorf("For input '%s' (%s): expected '%s', got '%s'",
					test.input, test.desc, test.expected, output)
			} else {
				t.Logf("Successfully reversliterated '%s' to '%s' (%s)",
					test.input, output, test.desc)
			}
		})
	}
}
