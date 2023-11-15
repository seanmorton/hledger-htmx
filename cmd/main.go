package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-htmx/internal/hledger"
	"github.com/seanmorton/hledger-htmx/internal/templates"
)

func main() {
	//accounts, _ := hledger.Accounts()
	register := []hledger.RegisterEntry{}
	balances, _ := hledger.Balances("li:cc")

	index := templates.Index(balances, register)
	http.Handle("/", templ.Handler(index))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		account := r.URL.Query().Get("account")
		register, _ := hledger.Register(account)
		templates.Register(register).Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", nil)
}
