package main

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/seanmorton/hledger-htmx/internal/hledger"
	"github.com/seanmorton/hledger-htmx/internal/templates"
)

//go:embed public
var publicDir embed.FS

//go:embed budget.json
var budgetContents []byte

// TODO error handling
func main() {
	budgetItems := []hledger.BudgetItem{}
	err := json.Unmarshal(budgetContents, &budgetItems)
	if err != nil {
		return
	}

	http.Handle("/public/", http.FileServer(http.FS(publicDir)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/budget", http.StatusFound)
	})

	http.HandleFunc("/budget", func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			from, to = defaultDateRange()
		}
		items, _ := hledger.Budget(from, to, budgetItems)
		render(w, r, templates.Budget(from, to, items))
	})

	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		acct := r.URL.Query().Get("account")
		var depth int
		if acct == "" || acct == "xp" || acct == "xp:" {
			acct = "xp:"
			depth = 2
		} else {
			depth = len(strings.Split(acct, ":")) + 1
		}
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			from, to = defaultDateRange()
		}
		balances, _ := hledger.Balances(acct, from, to, depth)
		register, _ := hledger.Register(acct, from, to)

		render(w, r, templates.Expenses(from, to, balances, register))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		acct := r.URL.Query().Get("account")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			from, to = defaultDateRange()
		}
		register, _ := hledger.Register(acct, from, to)

		render(w, r, templates.Register(from, to, register))
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
		from = time.Date(now.Year(), now.Month()-1, 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
		to = time.Date(now.Year(), now.Month(), 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	} else {
		from = time.Date(now.Year(), now.Month(), 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
		to = time.Date(now.Year(), now.Month()+1, 16, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	}
	return from, to
}
