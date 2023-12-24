package main

import (
	"embed"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-htmx/internal/hledger"
	"github.com/seanmorton/hledger-htmx/internal/templates"
)

//go:embed css
var FS embed.FS

func main() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.FS(FS))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/budget", http.StatusFound)
	})

	http.HandleFunc("/budget", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, templates.Budget())
	})

	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		account := r.URL.Query().Get("account")
		if account == "" {
			account = "xp"
		}
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			from, to = defaultDateRange()
		}
		balances, _ := hledger.Balances(account, from, to)
		register, _ := hledger.Register(account, from, to)

		render(w, r, templates.Expenses(account, from, to, balances, register))
	})

	http.ListenAndServe(":8080", nil)
}

func render(w http.ResponseWriter, r *http.Request, content templ.Component) {
	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
	} else {
		templates.Index(content).Render(r.Context(), w)
	}
}

func defaultDateRange() (string, string) {
	var from, to string
	now := time.Now()
	if now.Day() < 16 {
		from = time.Date(now.Year(), now.Month()-1, 15, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
		to = time.Date(now.Year(), now.Month(), 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	} else {
		from = time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
		to = time.Date(now.Year(), now.Month()+1, 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	}
	return from, to
}
