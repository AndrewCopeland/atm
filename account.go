package atm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Account struct {
	AccountID int
	PIN       string
	Balance   float64
}

type IAccountDB interface {
	// Return error if account cannot be found
	Get(int) (Account, error)
	// Return err if account could not be updated
	Set(Account) error
}

type AccountDB struct {
	DBFile string
}

func (a AccountDB) read() ([]Account, error) {
	file, err := os.Open(a.DBFile)
	if err != nil {
		return []Account{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	accounts := []Account{}
	for scanner.Scan() {
		// Skip first line since it is a CSV
		if firstLine {
			firstLine = false
			continue
		}
		line := scanner.Text()
		// Skip empty lines
		if line == "" {
			continue
		}

		columns := strings.Split(line, ",")

		// Validate all of the csv entries are appropriate types
		accountID, err := strconv.Atoi(columns[0])
		if err != nil {
			return []Account{}, ErrAccountIDNotInteger
		}

		pin := columns[1]

		balance, err := strconv.ParseFloat(columns[2], 64)
		if err != nil {
			return []Account{}, ErrAccountBalanceNotFloat
		}

		account := Account{
			AccountID: accountID,
			PIN:       pin,
			Balance:   balance,
		}
		accounts = append(accounts, account)
	}

	if err := scanner.Err(); err != nil {
		return []Account{}, err
	}

	return accounts, nil
}

func (a AccountDB) write(accounts []Account) error {
	file, err := os.OpenFile(a.DBFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)
	// Write the header
	datawriter.WriteString("ACCOUNT_ID,PIN,BALANCE\n")

	// Write each account to the file
	for _, account := range accounts {
		_, err = datawriter.WriteString(fmt.Sprintf("%d,%s,%.2f\n", account.AccountID, account.PIN, account.Balance))
		if err != nil {
			return err
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}

// Get returns a specific account from the CSV file
// if account cannot be found then an error is returned
func (a AccountDB) Get(accountID int) (Account, error) {
	accounts, err := a.read()
	if err != nil {
		return Account{}, err
	}
	for _, account := range accounts {
		if account.AccountID == accountID {
			return account, nil
		}
	}

	return Account{}, ErrAccountNotFound
}

// Set returns an error if the account was not updated in the CSV file
// set will override the CSV row that represents this account
func (a AccountDB) Set(account Account) error {
	accounts, err := a.read()
	if err != nil {
		return err
	}

	found := false
	for i, a := range accounts {
		if a.AccountID == account.AccountID {
			accounts[i] = account
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to update account because it does not exist")
	}

	err = a.write(accounts)
	if err != nil {
		return fmt.Errorf("Failed to write account to database. %s", err.Error())
	}
	return nil
}
