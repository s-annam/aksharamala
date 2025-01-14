// Package main implements the Aksharamala application, which provides functionalities
// for transliteration based on keymaps. It serves as the entry point for the application.
package main

import (
	"flag"

	"aks.go/internal/keymap"
	"aks.go/internal/translit"
	"aks.go/logger"
	"go.uber.org/zap"
)

// main is the entry point of the Aksharamala application.
// It parses command-line flags for configuration, initializes logging, and starts the application.
func main() {
	// Parse flags
	keymapsPath := flag.String("keymaps", "./keymaps", "Path to the keymaps directory")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Initialize the logger
	logger.InitLogger(*debug)
	defer logger.Sync()

	store := keymap.NewKeymapStore()
	if err := store.LoadKeymaps(*keymapsPath); err != nil {
		logger.Error("Failed to load keymaps", zap.String("path", *keymapsPath), zap.Error(err))
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
			logger.Error("Error during transliteration", zap.String("id", test.id), zap.String("input", test.input), zap.Error(err))
		} else {
			logger.Info("Transliteration successful", zap.String("input", test.input), zap.String("output", output))
		}
	}
}
