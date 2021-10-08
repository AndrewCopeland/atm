package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AndrewCopeland/atm"
)

var accountContent = []byte("ACCOUNT_ID,PIN,BALANCE\n12345678,1234,100.12")
var transactionContent = []byte("ACCOUNT_ID,DATE_TIME,AMOUNT,BALANCE\n")
var accountPath = "./accounts-e2e.csv"
var transactionPath = "./transactions-e2e.csv"

func writeTestContent(t *testing.T, path string, content []byte) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		t.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.Write(content) // Write at 0 beginning
	if err != nil {
		t.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func TestMainE2E(t *testing.T) {
	writeTestContent(t, accountPath, accountContent)
	writeTestContent(t, transactionPath, transactionContent)

	accountDB := atm.AccountDB{
		DBFile: accountPath,
	}
	transactionDB := atm.TransactionDB{
		DBFile: transactionPath,
	}
	a := &atm.ATM{
		AccountDB:     accountDB,
		TransactionDB: transactionDB,
		ATMBalance:    10000.00,
		Session:       &atm.Session{},
	}

	accountID := 12345678

	err := atm.RunCommand(a, "authorize 12345678 1234")
	if err != nil {
		t.Errorf("Failed to authorize")
	}

	err = atm.RunCommand(a, "withdraw 20")
	if err != nil {
		t.Errorf("Failed to withdraw 20")
	}

	err = atm.RunCommand(a, "balance")
	if err != nil {
		t.Errorf("Failed to get balance")
	}

	balance, _ := a.Balance(accountID)
	if balance != 80.12 {
		t.Errorf("Balance is invalid and should be 80.12")
	}

	err = atm.RunCommand(a, "deposit 20")
	if err != nil {
		t.Errorf("Failed to deposit money")
	}

	balance, _ = a.Balance(accountID)
	if balance != 100.12 {
		t.Errorf("Balance is invalid and should be 100.12")
	}

	err = atm.RunCommand(a, "history")
	if err != nil {
		t.Errorf("Failed to get history")
	}

	transactions := a.History(accountID)
	if len(transactions) != 2 {
		t.Errorf("Invalid number of transactions returned")
	}

	err = atm.RunCommand(a, "logout")
	if err != nil {
		t.Errorf("Failed to logout")
	}

	err = atm.RunCommand(a, "balance")
	if err == nil {
		t.Errorf("Error was expected when executing balanace when logged out")
	}

	os.Remove(accountPath)
	os.Remove(transactionPath)
}
