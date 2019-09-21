package main

import (
	"fmt"
	"io"
	"net/http"
)

type speechProcessor interface {
	ProcessText(text io.Reader) error
	GenerateRandomText(maxWords uint) string
}

// Handlers contains all the shared resources accross handlers
type Handlers struct {
	processor speechProcessor
	maxWords  uint
}

const textPlainContentType = "text/plain"

// Learn will parse the body of text on the request and process the ngrams
func (h *Handlers) Learn(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// If it's not text plain, don't process it
	if req.Header.Get("Content-Type") != textPlainContentType {
		http.Error(w, "Invalid content type", http.StatusUnprocessableEntity)
		return
	}

	if err := h.processor.ProcessText(req.Body); err != nil {
		http.Error(w, fmt.Sprintf("Error processing text: %v", err), http.StatusUnprocessableEntity)
		return
	}
}

// Generate will return randomly-generated text based on all the ngrams that
// have been learned since starting the program.
func (h *Handlers) Generate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var text = h.processor.GenerateRandomText(h.maxWords)

	w.Write([]byte(text))
}
