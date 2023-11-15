package hledger

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// TODO use correct struct types (or not..if R/O)
// TODO remove quotes from strings
// TODO Make stuct service for dep injection
// TODO properly escape csv, add comma back to commodity
// TODO splice slices for csv headers
// TODO Add start/end data params

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

func Balances(acct string) ([]BalanceEntry, error) {
	depth := strings.Count(acct, ":") + 2
	args := fmt.Sprintf("bal %s -%d -O csv", acct, depth)
	csvOutput, err := hledger(args)
	if err != nil {
		return nil, err
	}
	return parseBalances(csvOutput), nil
}

func Register(acct string) ([]RegisterEntry, error) {
	args := fmt.Sprintf("register %s -b 2023-09-16 -O csv", acct)
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

func parseBalances(csv string) []BalanceEntry {
	entries := []BalanceEntry{}
	rows := strings.Split(csv, "\n")

	for _, row := range rows {
		data := strings.Split(row, ",")
		if len(data) == 2 && data[0] != "\"account\"" && data[0] != "\"total\"" {
			amount, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(data[1], "\"", ""), "$", ""), 64)
			if err != nil {
				fmt.Println(err)
				return []BalanceEntry{}
			}
			entry := BalanceEntry{
				Account: strings.ReplaceAll(data[0], "\"", ""),
				Amount:  amount,
			}
			entries = append(entries, entry)
		}
	}
	return entries
}

func parseRegister(csv string) []RegisterEntry {
	entries := []RegisterEntry{}
	rows := strings.Split(csv, "\n")

	for _, row := range rows {
		data := strings.Split(row, ",")
		if len(data) == 7 && data[0] != "\"txnidx\"" {
			entry := RegisterEntry{
				Amount:      data[5],
				Account:     data[4],
				Date:        data[1],
				Description: data[3],
				Total:       data[6],
			}
			entries = append(entries, entry)
		}
	}
	return entries
}
