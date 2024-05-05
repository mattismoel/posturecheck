package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Embedding af templates: HTML-filer i "template"-mappen.
//
//go:embed template
var tmplFS embed.FS

// Embedding af statiske filer: CSS-filer i "static"-mappen.
//
//go:embed static
var staticFS embed.FS

// Antal gange at en "backpain" er registreret siden server start.
var checkCount = 0

func main() {
	mux := http.NewServeMux()

	// Gør statiske filer tilgængelige på "localhost:8080/static/...".
	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Registrering af stier på webserveren.
	mux.HandleFunc("/", handleIndex())
	mux.HandleFunc("POST /add", handleAdd())
	mux.Handle("GET /count", handleGetCount())

	log.Println("serving on :8080")

	// Opstart af serveren.
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Håndterer registrering af backpains ved POST request.
func handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkCount++
	}
}

// Håndterer rendering af startsiden.
func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gennemtjek nødvendige filer.
		// Hvis der er syntaksfejl eller lignende, skriver vi fejl til browseren.
		tmpl, err := template.ParseFS(
			tmplFS,
			"template/index.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Opsætnign af nødvendig data-map til vores startside.
		data := map[string]any{
			"checkCount": checkCount,
			"message":    "Some message...",
		}

		// Skriv respons til ResponseWriter (w), med ovenstående data.
		// Hvis fejl opstår, skriver vi fejl til browseren.
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Håndtering af GET request af ny counter.
func handleGetCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gennemtjek counter.html for eventuelle syntaksfejl.
		tmpl, err := template.ParseFS(tmplFS, "template/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Opsætning af nødvendig data for counter.html.
		data := map[string]any{
			"checkCount": checkCount,
			"message":    "Some other message...",
		}

		// Forsøg at skriv respons til request med den definerede data.
		err = tmpl.ExecuteTemplate(w, "content", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
