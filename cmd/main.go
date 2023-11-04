package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-webapp/internal"
	"github.com/seanmorton/hledger-webapp/internal/templates"
)

func main() {
	accounts, _ := internal.Accounts()
	register := []internal.RegisterEntry{}
	//balances, _ := internal.Balances("as:stocks")

	index := templates.Index(accounts, register)
	http.Handle("/", templ.Handler(index))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		account := r.URL.Query().Get("account")
		register, _ := internal.Register(account)
		templates.Register(register).Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", nil)
}
