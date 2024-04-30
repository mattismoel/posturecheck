package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", handleAdd())
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Added one"))
	}
}
