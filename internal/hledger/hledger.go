package hledger

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type Balance struct {
	Account     string    `json:"account"`
	Amount      float64   `json:"amount"`
	SubBalances []Balance `json:"sub_balances"`
}

type RegisterEntry struct {
	Account     string
	Amount      float64
	Date        string
	Description string
	Total       float64
}

type BudgetItem struct {
	Name    string  `json:"name"`
	Account string  `json:"account"`
	Target  float64 `json:"target"`
	Spent   float64 `json:"spent"`
}

func (b *BudgetItem) Remaining() float64 {
	return b.Target - b.Spent
}

func (b *BudgetItem) Percent() float64 {
	return b.Remaining() / b.Target * 100
}

func Balances(acct, from, to string, depth int, historical, invert bool) (Balance, error) {
	args := fmt.Sprintf("bal %s -b %s -e %s -S -V -O csv", acct, from, to)
	if depth > 0 {
		args = fmt.Sprintf("%s -%d", args, depth)
	}
	if historical {
		args = fmt.Sprintf("%s --historical", args)
	}
	if invert {
		args = fmt.Sprintf("%s --invert", args)
	}
	csvOutput, err := hledger(args)
	if err != nil {
		return Balance{}, err
	}
	return parseBalances(acct, csvOutput), nil
}

func Register(acct, to, from string, historical, invert bool) ([]RegisterEntry, error) {
	args := fmt.Sprintf("register %s -b %s -e %s -V -O csv", acct, to, from)
	if historical {
		args = fmt.Sprintf("%s --historical", args)
	}
	if invert {
		args = fmt.Sprintf("%s --invert", args)
	}
	csvOutput, err := hledger(args)
	if err != nil {
		return nil, err
	}
	return parseRegister(csvOutput), nil
}

func Budget(from, to string, items []BudgetItem) ([]BudgetItem, error) {
	for i, item := range items {
		balance, _ := Balances(item.Account, from, to, 0, false, false)
		item.Spent = balance.Amount
		items[i] = item
	}
	return items, nil
}

// TODO(security) prevent malicious command injection
func hledger(args string) (string, error) {
	cmd := exec.Command("hledger", strings.Split(args, " ")...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(out)
	return buf.String(), nil
}

func parseBalances(acct, csv string) Balance {
	var total float64
	subBalances := []Balance{}
	rows := strings.Split(csv, "\n")

	for _, row := range rows {
		data := strings.Split(row, ",")
		if len(data) == 2 && data[0] != "\"account\"" { // Skip header
			amount, _ := parseAmount(data[1])
			if data[0] == "\"Total:\"" {
				total = amount
				continue
			}
			balance := Balance{
				Account: strings.ReplaceAll(data[0], "\"", ""),
				Amount:  amount,
			}
			subBalances = append(subBalances, balance)
		}
	}

	if strings.HasSuffix(acct, ":") {
		acct = acct[:len(acct)-1]
	}
	return Balance{
		Account:     acct,
		Amount:      total,
		SubBalances: subBalances,
	}
}

func parseRegister(csv string) []RegisterEntry {
	entries := []RegisterEntry{}
	rows := strings.Split(csv, "\n")

	for _, row := range rows {
		data := strings.Split(row, ",")
		if len(data) == 7 && data[0] != "\"txnidx\"" {
			amount, _ := parseAmount(data[5])
			total, _ := parseAmount(data[6])
			entry := RegisterEntry{
				Amount:      amount,
				Account:     strings.ReplaceAll(data[4], "\"", ""),
				Date:        strings.ReplaceAll(data[1], "\"", ""),
				Description: strings.ReplaceAll(data[3], "\"", ""),
				Total:       total,
			}
			entries = append(entries, entry)
		}
	}
	slices.Reverse(entries) // Show most recent first
	return entries
}

func parseAmount(amount string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(amount, "\"", ""), "$", ""), 64)
}
