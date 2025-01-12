package main

import (
	"fmt"
	"log"

	"aks.go/internal/keymap"
	"aks.go/internal/translit"
)

func main() {
	// Initialize keymap store
	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps("./keymaps"); err != nil {
		log.Fatalf("Failed to load keymaps: %v", err)
	}

	// Get Hindi scheme
	scheme, exists := store.GetKeymap("hindi")
	if !exists {
		log.Fatal("Hindi keymap not found")
	}

	// Create Aksharamala instance
	aks, err := translit.NewAksharamala(&scheme)
	if err != nil {
		log.Fatalf("Failed to create Aksharamala: %v", err)
	}

	// Example transliteration
	input := "namaste"
	output, err := aks.Transliterate(input)
	if err != nil {
		log.Fatalf("Transliteration failed: %v", err)
	}

	fmt.Printf("Input text: %s\n", input)
	fmt.Printf("Transliterated text: %s\n", output)
}
