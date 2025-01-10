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
