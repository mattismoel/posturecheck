package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed template
var tmplFS embed.FS

var CheckCount = 0

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", handleAdd())
	mux.HandleFunc("/", handleIndex())

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		CheckCount++
		w.Write([]byte("Added one"))
	}
}

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(tmplFS, "template/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, CheckCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
