package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed template
var tmplFS embed.FS

//go:embed static
var staticFS embed.FS

var checkCount = 0

func main() {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("POST /add", handleAdd())
	mux.Handle("GET /count", handleGetCount())
	mux.HandleFunc("/", handleIndex())

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkCount++
		w.Write([]byte("Added one"))
	}
}

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(
			tmplFS,
			"template/index.html",
			"template/counter.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"checkCount": checkCount,
			"message":    "Some message...",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGetCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(tmplFS, "template/counter.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"checkCount": checkCount,
			"message":    "Some other message...",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
