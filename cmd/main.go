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
	defaultAccount := "xp"
	defaultFrom, defaultTo := defaultDateRange()
	balances, _ := hledger.Balances(defaultAccount, defaultFrom, defaultTo)
	register, _ := hledger.Register(defaultAccount, defaultFrom, defaultTo)

	index := templates.Index(defaultAccount, defaultFrom, defaultTo, balances, register)
	http.Handle("/", templ.Handler(index))

	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		account := r.URL.Query().Get("account")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		balances, _ := hledger.Balances(account, from, to)
		register, _ := hledger.Register(account, from, to)

		templates.Expenses(account, from, to, balances, register).Render(r.Context(), w)
	})

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.FS(FS))))

	http.ListenAndServe(":8080", nil)
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
