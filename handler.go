package main

import (
	"html/template"
	"net/http"
)

// Håndterer registrering af backpains ved POST request.
func handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkCount++
	}
}

// Håndterer rendering af startsiden.
func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieCount, err := cookieCount(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		checkCount = cookieCount

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
		setCookieCount(w, checkCount)

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
