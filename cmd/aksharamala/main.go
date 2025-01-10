package main

import (
	"fmt"

	"aks.go/internal/keymap"
)

func main() {
	store := keymap.NewKeymapStore()

	// Load keymaps from the "keymaps" directory
	err := store.LoadKeymaps("./keymaps")
	if err != nil {
		fmt.Printf("Error loading keymaps: %v\n", err)
		return
	}

	// List all loaded keymaps
	fmt.Println("Loaded Keymaps:")
	for _, id := range store.ListKeymapIDs() {
		fmt.Println("- ", id)
	}
}
