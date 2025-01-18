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
			{ID: "Hindi", Name: "Hindi (ITRANS -> Unicode)"},
			{ID: "Marathi", Name: "Marathi (ITRANS -> Unicode)"},
			{ID: "TeluguRts", Name: "Telugu (RTS -> Unicode)"},
			{ID: "RHindi", Name: "Hindi (Unicode -> ITRANS)"},
			{ID: "RSanskrit", Name: "Sanskrit (Unicode -> ITRANS)"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keymaps)
	})

	// Map the given text using the selected keymap
	http.HandleFunc("/api/m", func(w http.ResponseWriter, r *http.Request) {
		// Handle OPTIONS (CORS preflight request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var text, keymapID string

		if r.Method == http.MethodGet {
			// Extract parameters from the URL for GET requests
			query := r.URL.Query()
			text = query.Get("text")
			keymapID = query.Get("keymapId")

			if text == "" || keymapID == "" {
				http.Error(w, "Missing required parameters: text and keymapId", http.StatusBadRequest)
				return
			}
		} else if r.Method == http.MethodPost {
			// Decode JSON request body for POST requests
			var req TransliterationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			text = req.Text
			keymapID = req.KeymapID
		} else {
			// Reject other methods
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Perform transliteration
		result, err := aksharamala.TransliterateWithKeymap(keymapID, text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := TransliterationResponse{Result: result}

		// ✅ Ensure correct UTF-8 encoding
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
