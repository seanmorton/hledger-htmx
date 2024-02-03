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
	Amount      string
	Date        string
	Description string
	Total       string
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

func Balances(acct string, from, to string, depth int) (Balance, error) {
	var args string
	if depth > 0 {
		args = fmt.Sprintf("bal %s -%d -b %s -e %s -S -O csv", acct, depth, from, to)
	} else {
		args = fmt.Sprintf("bal %s -b %s -e %s -S -O csv", acct, from, to)
	}
	csvOutput, err := hledger(args)
	if err != nil {
		return Balance{}, err
	}
	return parseBalances(acct, csvOutput), nil
}

func Register(acct string, to, from string) ([]RegisterEntry, error) {
	args := fmt.Sprintf("register %s -b %s -e %s -O csv", acct, to, from)
	csvOutput, err := hledger(args)
	if err != nil {
		return nil, err
	}
	return parseRegister(csvOutput), nil
}

func Budget(from, to string, items []BudgetItem) ([]BudgetItem, error) {
	for i, item := range items {
		balance, _ := Balances(item.Account, from, to, 0)
		item.Spent = balance.Amount
		items[i] = item
	}
	return items, nil
}

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
			amount, _ := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(data[1], "\"", ""), "$", ""), 64)
			if data[0] == "\"total\"" {
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
			entry := RegisterEntry{
				Amount:      strings.ReplaceAll(data[5], "\"", ""),
				Account:     strings.ReplaceAll(data[4], "\"", ""),
				Date:        strings.ReplaceAll(data[1], "\"", ""),
				Description: strings.ReplaceAll(data[3], "\"", ""),
				Total:       strings.ReplaceAll(data[6], "\"", ""),
			}
			entries = append(entries, entry)
		}
	}
	slices.Reverse(entries)
	return entries
}
