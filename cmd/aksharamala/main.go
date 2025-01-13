package main

import (
	"fmt"

	"aks.go/internal/translit"
)

// Example test setup
func main() {
	schemePath := "c:\\repo\\aks.go\\keymaps\\Hindi.aksj"
	aks, err := translit.NewAksharamala(schemePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	inputs := []string{"namaste", "kk", "ka"}
	for _, input := range inputs {
		output := aks.Transliterate(input)
		fmt.Printf("Input: %s\nOutput: %s\n", input, output)
	}
}
