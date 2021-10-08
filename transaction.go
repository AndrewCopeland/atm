package atm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {
	AccountID int
	DateTime  int64
	Amount    float64
	Balance   float64
}

type ITransactionDB interface {
	// return all transactions for a given account
	Get(int) ([]Transaction, error)
	// add a transaction return error if failure to add transaction
	Set(Transaction) error
}

type TransactionDB struct {
	DBFile string
}

func (t TransactionDB) read() ([]Transaction, error) {
	file, err := os.Open(t.DBFile)
	if err != nil {
		return []Transaction{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	transactions := []Transaction{}
	for scanner.Scan() {
		// Skip first line since it is a CSV
		if firstLine {
			firstLine = false
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}

		columns := strings.Split(line, ",")

		// Validate all of the csv entries of are appropriate types
		accountID, err := strconv.Atoi(columns[0])
		if err != nil {
			return []Transaction{}, ErrTransactionAccountIDNotInteger
		}

		dateTime, err := strconv.ParseInt(columns[1], 10, 64)
		if err != nil {
			return []Transaction{}, ErrTransactionDateTimeNotInteger
		}

		amount, err := strconv.ParseFloat(columns[2], 64)
		if err != nil {
			return []Transaction{}, ErrTransactionAmounteNotFloat
		}

		balance, err := strconv.ParseFloat(columns[3], 64)
		if err != nil {
			return []Transaction{}, ErrTransactionBalanceNotFloat
		}

		transaction := Transaction{
			AccountID: accountID,
			DateTime:  dateTime,
			Amount:    amount,
			Balance:   balance,
		}
		transactions = append(transactions, transaction)
	}

	if err := scanner.Err(); err != nil {
		return []Transaction{}, err
	}

	return transactions, nil
}

func (t TransactionDB) write(transactions []Transaction) error {
	file, err := os.OpenFile(t.DBFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)
	// Write the header
	datawriter.WriteString("ACCOUNT_ID,DATE_TIME,AMOUNT,BALANCE\n")

	// Write each transaction to file
	for _, transaction := range transactions {
		_, err = datawriter.WriteString(fmt.Sprintf("%d,%d,%.2f,%.2f\n", transaction.AccountID, transaction.DateTime, transaction.Amount, transaction.Balance))

		if err != nil {
			return err
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}

// Get retrieves a list of transactions for a given accountID from a CSV file
// An error is returned on failure to read the CSV file
func (t TransactionDB) Get(accountID int) ([]Transaction, error) {
	transactions, err := t.read()
	if err != nil {
		return []Transaction{}, err
	}

	accountTransactions := []Transaction{}
	for _, transaction := range transactions {
		if transaction.AccountID == accountID {
			accountTransactions = append(accountTransactions, transaction)
		}
	}

	return accountTransactions, err
}

// Set appends a transaction to the transactions CSV file
// An error is returned on failure to read or write the CSV file
func (t TransactionDB) Set(transaction Transaction) error {
	transactions, err := t.read()
	if err != nil {
		return err
	}

	transactions = append(transactions, transaction)
	err = t.write(transactions)
	if err != nil {
		return err
	}
	return nil
}
