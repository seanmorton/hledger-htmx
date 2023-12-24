package hledger

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// TODO use correct struct types (or not..if R/O)
// TODO properly escape csv, add comma back to commodity
// TODO splice slices for csv headers

type Balance struct {
	Account string         `json:"account"`
	Total   float64        `json:"total"`
	Entries []BalanceEntry `json:"entries"`
}

type BalanceEntry struct {
	Account string  `json:"account"`
	Amount  float64 `json:"amount"`
}

type RegisterEntry struct {
	Account     string
	Amount      string
	Date        string
	Description string
	Total       string
}

func Accounts() ([]string, error) {
	output, err := hledger("accounts")
	if err != nil {
		return nil, err
	}
	return parseAccounts(output), nil
}

func Balances(acct string, to, from string, depth int) (Balance, error) {
	// If showing a root account, add a colon to avoid pulling in unintended accounts
	if strings.Count(acct, ":") == 0 {
		acct += ":"
	}
	args := fmt.Sprintf("bal %s -%d -b %s -e %s -S -O csv", acct, depth, to, from)
	csvOutput, err := hledger(args)
	if err != nil {
		return Balance{}, err
	}
	return parseBalances(acct, csvOutput), nil
}

func Register(acct string, to, from string) ([]RegisterEntry, error) {
	// If showing a root account, add a colon to avoid pulling in unintended accounts
	if strings.Count(acct, ":") == 0 {
		acct += ":"
	}
	args := fmt.Sprintf("register %s -b %s -e %s -O csv", acct, to, from)
	csvOutput, err := hledger(args)
	if err != nil {
		return nil, err
	}
	return parseRegister(csvOutput), nil
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

func parseAccounts(output string) []string {
	rows := strings.Split(output, "\n")
	return rows
}

func parseBalances(acct, csv string) Balance {
	var total float64
	entries := []BalanceEntry{}
	rows := strings.Split(csv, "\n")

	for _, row := range rows {
		data := strings.Split(row, ",")
		if len(data) == 2 && data[0] != "\"account\"" { // Skip header
			amount, _ := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(data[1], "\"", ""), "$", ""), 64)
			if data[0] == "\"total\"" {
				total = amount
				continue
			}
			entry := BalanceEntry{
				Account: strings.ReplaceAll(data[0], "\"", ""),
				Amount:  amount,
			}
			entries = append(entries, entry)
		}
	}

	return Balance{
		Account: acct,
		Total:   total,
		Entries: entries,
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
