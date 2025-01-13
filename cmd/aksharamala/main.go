package main

import (
	"fmt"

	"aks.go/internal/translit"
)

// Example test setup
func main() {
	schemePath := "keymaps\\Hindi.aksj"
	aks, err := translit.NewAksharamala(schemePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	inputs := []string{"namaste", "a1k"}
	for _, input := range inputs {
		output := aks.Transliterate(input)
		fmt.Printf("Input: %s\nOutput: %s\n", input, output)
	}
}
