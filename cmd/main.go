package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-htmx/internal/hledger"
	"github.com/seanmorton/hledger-htmx/internal/templates"
)

func main() {
	//accounts, _ := hledger.Accounts()
	balances, _ := hledger.Balances("xp:primary")
	register, _ := hledger.Register("xp:primary")

	index := templates.Index(balances, register)
	http.Handle("/", templ.Handler(index))

	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		account := r.URL.Query().Get("account")
		balances, _ := hledger.Balances(account)
		register, _ := hledger.Register(account)
		templates.Expenses(balances, register).Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", nil)
}
