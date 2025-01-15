package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"aks.go/internal/keymap"
	"aks.go/internal/translit"
)

type TransliterationRequest struct {
	Text     string `json:"text"`
	KeymapID string `json:"keymapId"`
}

type TransliterationResponse struct {
	Result string `json:"result"`
}

type Keymap struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var aksharamala *translit.Aksharamala

func init() {
	// Initialize keymap store with keymaps from the keymaps directory
	store := keymap.NewKeymapStore()
	keymapsDir, err := filepath.Abs("./keymaps")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if err := store.LoadKeymaps(keymapsDir); err != nil {
		log.Printf("Warning: Failed to load keymaps: %v", err)
	}

	aksharamala = translit.NewAksharamala(store)
}

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get available keymaps
	http.HandleFunc("/api/keymaps", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get keymaps from the store
		keymaps := []Keymap{
			{ID: "Hindi", Name: "Hindi (Devanagari)"},
			{ID: "Marathi", Name: "Marathi (Devanagari)"},
			{ID: "TeluguRts", Name: "Telugu (RTS)"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keymaps)
	})

	// Transliterate endpoint
	http.HandleFunc("/api/transliterate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req TransliterationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Perform transliteration using your engine
		result, err := aksharamala.TransliterateWithKeymap(req.KeymapID, req.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := TransliterationResponse{
			Result: result,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Enable CORS
	handler := corsMiddleware(http.DefaultServeMux)

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
