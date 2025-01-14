// This file is part of Aksharamala (aks.go).
//
// Aksharamala is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Aksharamala is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Aksharamala. If not, see <https://www.gnu.org/licenses/>.

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
// It provides methods to load keymaps from JSON files, retrieve keymaps by ID,
// and list all loaded keymap IDs. The store is thread-safe due to the use
// of a read-write mutex.
type KeymapStore struct {
	// Maps keymap IDs to TransliterationScheme
	Keymaps map[string]types.TransliterationScheme
	// Mutex for concurrent access
	mu sync.RWMutex
}

// NewKeymapStore initializes a new KeymapStore with an empty map of keymaps.
func NewKeymapStore() *KeymapStore {
	return &KeymapStore{
		Keymaps: make(map[string]types.TransliterationScheme),
	}
}

// LoadKeymaps loads JSON keymaps from a specified directory into the store.
// It reads all JSON files in the directory and adds them to the Keymaps map.
// Returns an error if loading any keymap fails.
func (store *KeymapStore) LoadKeymaps(directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".aksj") {
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
// It reads the specified file and updates the Keymaps map with its contents.
// Returns an error if the file cannot be read or if the JSON is invalid.
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

	// Validate the scheme
	if err := scheme.Validate(); err != nil {
		return fmt.Errorf("keymap validation failed for '%s': %w", filePath, err)
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	store.Keymaps[scheme.ID] = scheme
	return nil
}

// GetKeymap retrieves a TransliterationScheme by ID.
// Returns the TransliterationScheme and a boolean indicating whether it was found.
func (store *KeymapStore) GetKeymap(id string) (types.TransliterationScheme, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	scheme, exists := store.Keymaps[id]
	return scheme, exists
}

// ListKeymapIDs returns a list of all loaded keymap IDs.
// This is useful for iterating over available keymaps.
func (store *KeymapStore) ListKeymapIDs() []string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	ids := make([]string, 0, len(store.Keymaps))
	for id := range store.Keymaps {
		ids = append(ids, id)
	}

	return ids
}
