package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
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

var countCookieName = "check-count"

func main() {
	mux := http.NewServeMux()

	// Gør statiske filer tilgængelige på "localhost:8080/static/...".
	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	timer := time.NewTicker(timeTillMidnight())
	go func() {
		for range timer.C {
			checkCount = 0
			timer.Reset(timeTillMidnight())
		}
	}()

	// Statiske filer håndteres på 'localhost:PORT/static/...'
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Registrering af stier på webserveren.
	mux.HandleFunc("/", handleIndex())
	mux.HandleFunc("POST /add", handleAdd())
	mux.Handle("GET /count", handleGetCount())

	log.Println("serving on :8080")

	// Opstart af serveren.
	log.Fatal(http.ListenAndServe(":8080", mux))
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

func cookieCount(r *http.Request) (int, error) {
	c, err := r.Cookie(countCookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return 0, nil
		default:
			return -1, fmt.Errorf("could not get cookie: %v", err)
		}
	}

	count, err := strconv.Atoi(c.Value)
	if err != nil {
		return -1, fmt.Errorf("could not convert cookie count to int: %v", err)
	}

	return count, nil
}

func setCookieCount(w http.ResponseWriter, count int) {
	c := &http.Cookie{
		Name:     countCookieName,
		Value:    strconv.Itoa(count),
		Path:     "/",
		Expires:  time.Now().Add(24 * 365 * 10 * time.Hour),
		HttpOnly: false,
		Secure:   false,
	}

	http.SetCookie(w, c)
}

func timeTillMidnight() time.Duration {
	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		log.Fatal(err)
	}

	// Bestem mængde af tid til midnat, og lav en timer (ticker), som ved 00:00
	// nulstiller 'checkCount'.
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	duration := midnight.Sub(now)
	return duration
}
