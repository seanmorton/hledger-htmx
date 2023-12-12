package main

import (
	"embed"
	"net/http"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-htmx/internal/hledger"
	"github.com/seanmorton/hledger-htmx/internal/templates"
)

//go:embed css
var FS embed.FS

func main() {
	defaultAccount := "xp"
	balances, _ := hledger.Balances(defaultAccount, "2023-11-15", "2023-11-30")
	register, _ := hledger.Register(defaultAccount, "2023-11-15", "2023-11-30")

	index := templates.Index(defaultAccount, balances, register)
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
