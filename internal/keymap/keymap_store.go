package keymap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"aks.go/internal/types"
)

// KeymapStore is an in-memory storage for all loaded keymaps.
type KeymapStore struct {
	Keymaps map[string]types.TransliterationScheme
	mu      sync.RWMutex
}

// NewKeymapStore initializes a new KeymapStore.
func NewKeymapStore() *KeymapStore {
	return &KeymapStore{
		Keymaps: make(map[string]types.TransliterationScheme),
	}
}

// LoadKeymaps loads JSON keymaps from a specified directory into the store.
func (store *KeymapStore) LoadKeymaps(directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		if err := store.loadKeymapFromFile(filePath); err != nil {
			return fmt.Errorf("failed to load keymap from file %s: %w", filePath, err)
		}
	}

	return nil
}

// loadKeymapFromFile loads a single JSON keymap file into the store.
func (store *KeymapStore) loadKeymapFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var scheme types.TransliterationScheme
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&scheme); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	if scheme.ID == "" {
		return fmt.Errorf("missing id in keymap file: %s", filePath)
	}

	// Safely add the keymap to the store
	store.mu.Lock()
	defer store.mu.Unlock()
	store.Keymaps[scheme.ID] = scheme

	return nil
}

// GetKeymap retrieves a TransliterationScheme by ID.
func (store *KeymapStore) GetKeymap(id string) (types.TransliterationScheme, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	scheme, exists := store.Keymaps[id]
	return scheme, exists
}

// ListKeymapIDs returns a list of all loaded keymap IDs.
func (store *KeymapStore) ListKeymapIDs() []string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	ids := make([]string, 0, len(store.Keymaps))
	for id := range store.Keymaps {
		ids = append(ids, id)
	}

	return ids
}
