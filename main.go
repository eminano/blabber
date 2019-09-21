package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/eminano/markov"
)

func main() {
	var n = flag.Uint("ngram", 3, "size of the ngrams to be processed when using learn")
	var maxWords = flag.Uint("maxWords", 100, "max number of words on output when using generate")
	var port = flag.Uint("port", 8080, "port for the server to listen on")

	flag.Parse()

	var ngramChain, err = markov.NewNGramChain(*n)
	if err != nil {
		log.Fatalf("error initialising text processor: %v", err)
	}

	h := Handlers{
		processor: ngramChain,
		maxWords:  *maxWords,
	}

	http.HandleFunc("/learn", h.Learn)
	http.HandleFunc("/generate", h.Generate)

	log.Printf("Listening on localhost:%d", *port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("error starting local server: %v", err)
	}
}
