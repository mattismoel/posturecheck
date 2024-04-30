package main

import (
	"fmt"
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
		fmt.Println("Added one")
		w.Write([]byte("Added one"))
	}
}
