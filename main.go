package main

import (
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
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

var db *sql.DB

func main() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS users (
    id            INTGER     PRIMARY KEY,
    username      TEXT      NOT NULL,
    check_count   INTEGER   NOT NULL DEFAULT 0
  )`)
	if err != nil {
		log.Fatal(err)
	}

	query := `SELECT check_count FROM users WHERE username = ?`
	err = db.QueryRow(query, "test_user").Scan(&checkCount)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Gør statiske filer tilgængelige på "localhost:8080/static/...".
	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	duration := midnight.Sub(now)
	timer := time.NewTicker(duration)

	go func() {
		<-timer.C
		checkCount = 0
	}()

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
		checkCount += 1
		query := `
    UPDATE users
      SET check_count = ?
    WHERE username = ?`

		_, err := db.Exec(query, checkCount, "test_user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

		msg := countToMsg(checkCount)

		// Opsætnign af nødvendig data-map til vores startside.
		data := map[string]any{
			"checkCount": checkCount,
			"message":    msg,
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

		msg := countToMsg(checkCount)

		// Opsætning af nødvendig data for counter.html.
		data := map[string]any{
			"checkCount": checkCount,
			"message":    msg,
		}

		// Forsøg at skriv respons til request med den definerede data.
		err = tmpl.ExecuteTemplate(w, "counter", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func countToMsg(count int) string {
	var message string // ""
	if count <= 10 {
		message = "wow good job today lad"
	}

	if count > 10 {
		message = "You're looking nasty at age 40..."
	}

	if count > 20 {
		message = "You're looking nasty at age 20..."
	}

	return message
}
