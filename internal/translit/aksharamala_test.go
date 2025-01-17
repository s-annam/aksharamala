package translit

import (
	"fmt"
	"testing"

	"aks.go/internal/keymap"
)

// TestTransliterate tests the Transliterate method of the Aksharamala struct.
// It verifies that the transliteration works correctly for various input cases.
// The test loads keymaps and checks the output against expected results.
func TestTransliterate(t *testing.T) {
	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps("../../keymaps"); err != nil {
		fmt.Printf("Failed to load keymaps: %v\n", err)
		return
	}

	aks := NewAksharamala(store)

	hindiSmallTests := []struct {
		input    string
		expected string
	}{
		{"namaste", "नमस्ते"},
		{"kk", "क्क"},
		{"ka", "क"},
		{"a1k", "अ१क"},
		{"", ""}, // Edge case: empty string
	}

	hindiLargeTests := []struct {
		input    string
		expected string
	}{
		{
			"yah ek su.ndar din hai. aaj ham bahut khush hai.n. hamane bahut saaree cheeze.n seekhee hai.n.",
			"यह एक सुंदर दिन है। आज हम बहुत खुश हैं। हमने बहुत सारी चीज़ें सीखी हैं।",
		},
		{
			"aur ab ham tayyaar hai.n naye anubhavo.n ke liye.",
			"और अब हम तय्यार हैं नये अनुभवों के लिये।",
		},
	}

	aks.SetActiveKeymap("hindi")

	// Consolidate all the Hindi tests and run them
	allTests := append(hindiSmallTests, hindiLargeTests...)
	for _, test := range allTests {
		translit(t, aks, test)
	}

	// Test with TeluguRTS keymap
	aks.SetActiveKeymap("teluguRts")

	teluguSmallTests := []struct {
		input    string
		expected string
	}{
		{"namastE", "నమస్తే"},
		{"kka", "క్క"},
		{"ka", "క"},
		{"a1k", "అ1క్"},
		{"daas", "దాస్"},
		{"jeevitam", "జీవితం"},
		{"avasaram.", "అవసరం."},
		{"", ""}, // Edge case: empty string
	}

	teluguLargeTests := []struct {
		input    string
		expected string
	}{
		{
			"ee rOju oka AnaMdamaina sudinamu. nEnu chAlA saMtOshaMgaa kotta vishayaalu nErchukOvaDaMlO utsaahaMgaa unnAnu.",
			"ఈ రోజు ఒక ఆనందమైన సుదినము. నేను చాలా సంతోషంగా కొత్త విషయాలు నేర్చుకోవడంలో ఉత్సాహంగా ఉన్నాను.",
			// ఈ రోజు ఒక ఆనందమైంఅ సుదినము. నేను చాలా సంతోషంగా కొత్త విషయాలు నేర్చుకోవడంలో ఉత్సాహంగా ఉన్నాను.
		},
		{
			"jeevitam aaScharyaala tO niMDinadi. prati avakaasAnni dhairyaMtO mariyu nammakaMtO sviikariMchaDam manaku avasaram.",
			"జీవితం ఆశ్చర్యాల తో నిండినది. ప్రతి అవకాసాన్ని ధైర్యంతో మరియు నమ్మకంతో స్వీకరించడం మనకు అవసరం.",
		},
	}

	// Consolidate all the Telugu RTS tests and run them
	allTests = append(teluguSmallTests, teluguLargeTests...)
	for _, test := range allTests {
		translit(t, aks, test)
	}
}

func TestContextualMapping(t *testing.T) {
	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps("../../keymaps"); err != nil {
		t.Fatalf("Failed to load keymaps: %v", err)
	}

	aks := NewAksharamala(store)
	aks.SetActiveKeymap("teluguRts")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple context update",
			input:    "kk",
			expected: "క్క్", // Assuming TeluguRts.aksj has the contextual mapping for 'k'
		},
		{
			name:     "No context match",
			input:    "k",
			expected: "క్", // No context, no modification
		},
		{
			name:     "Multiple context changes",
			input:    "kkk",
			expected: "క్క్క్", // Testing chain of contextual rules
		},
		{
			name:     "Context with space",
			input:    "kk ",
			expected: "క్క్ ", // Context should persist across space
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "", // Edge case
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := aks.Transliterate(test.input)
			if result != test.expected {
				t.Errorf("Test %s failed:\nexpected: %q\ngot: %q",
					test.name, test.expected, result)
			}
		})
	}
}

func translit(t *testing.T, aks *Aksharamala, test struct {
	input    string
	expected string
}) {
	output := aks.Transliterate(test.input)
	if output != test.expected {
		t.Errorf("For input '%s': expected '%s', got '%s'", test.input, test.expected, output)
	} else {
		t.Logf("For input '%s': output matches expected '%s'", test.input, test.expected)
	}
}
