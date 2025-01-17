package translit

import (
	"fmt"

	"aks.go/internal/keymap"
	"aks.go/internal/types"
)

// Aksharamala represents a transliteration engine that uses a keymap store
// to perform transliteration operations. It maintains the active transliteration
// scheme and manages the context for transliteration processes.
type Aksharamala struct {
	keymapStore   *keymap.KeymapStore
	activeScheme  *types.TransliterationScheme
	context       *types.Context
	viramaHandler *types.ViramaHandler
}

// NewAksharamala initializes a new Aksharamala instance.
func NewAksharamala(store *keymap.KeymapStore) *Aksharamala {
	return &Aksharamala{
		keymapStore: store,
		context:     types.NewContext(),
	}
}

// SetActiveKeymap sets the active keymap by ID for transliteration.
func (a *Aksharamala) SetActiveKeymap(id string) error {
	scheme, exists := a.keymapStore.GetKeymap(id)
	if !exists {
		return fmt.Errorf("keymap with ID '%s' not found", id)
	}

	virama, viramaMode, err := types.ParseVirama(scheme.Metadata.Virama)
	if err != nil {
		return fmt.Errorf("failed to parse virama: %v", err)
	}

	a.activeScheme = &scheme
	a.context = types.NewContext()
	a.viramaHandler = types.NewViramaHandler(viramaMode, virama, a.context)
	return nil
}

// TransliterateWithKeymap performs mapping of the input string using the active keymap.
func (a *Aksharamala) TransliterateWithKeymap(id, input string) (string, error) {
	if err := a.SetActiveKeymap(id); err != nil {
		return "", err
	}
	switch a.activeScheme.Scheme {
	case "Unicode":
		return a.Reversliterate(input)
	case "ITRANS", "RTS":
		return a.Transliterate(input)
	}
	return "", fmt.Errorf("unsupported scheme: %s", a.activeScheme.Scheme)
}
