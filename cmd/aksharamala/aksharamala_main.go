package main

import (
	"flag"
	"fmt"

	"aks.go/internal/keymap"
	"aks.go/internal/translit"
)

// Example test setup
func main() {
	// Define a flag for the keymaps directory
	keymapsPath := flag.String("keymaps", "./keymaps", "Path to the keymaps directory")
	flag.Parse()

	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps(*keymapsPath); err != nil {
		fmt.Printf("Failed to load keymaps: %v\n", err)
		return
	}

	aks := translit.NewAksharamala(store)

	inputs := []struct {
		id    string
		input string
	}{
		{"hindi", "namaste"},
		{"hindi", "ka"},
	}

	for _, test := range inputs {
		output, err := aks.TransliterateWithKeymap(test.id, test.input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Input: %s\nOutput: %s\n", test.input, output)
		}
	}
}
