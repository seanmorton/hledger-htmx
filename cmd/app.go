package main

import (
	"fmt"

	"github.com/seanmorton/hledger-api/internal"
)

func main() {
	balances, err := internal.Balances("as:stocks")
	if err != nil {
		panic(err)
	}
	for _, balance := range balances {
		fmt.Printf("%v\n", balance)
	}
	fmt.Println()

	register, err := internal.Register("as:cash:wealthfront")
	if err != nil {
		panic(err)
	}
	for _, entry := range register {
		fmt.Println(entry)
	}

	accounts, err := internal.Accounts()
	if err != nil {
		panic(err)
	}
	for _, account := range accounts {
		fmt.Printf("%v\n", account)
	}
	fmt.Println()
}
