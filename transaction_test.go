package atm_test

import (
	"testing"
	"time"

	"github.com/AndrewCopeland/atm"
)

var TransactionDB = atm.TransactionDB{
	DBFile: "./transactions_test.csv",
}

var InvalidTransactionDB = atm.TransactionDB{
	DBFile: "./transactions_not_real.csv",
}

func TestGetValidTransactions(t *testing.T) {
	transactions, err := TransactionDB.Get(7089382418)
	assertNoError(t, err)

	if len(transactions) != 3 {
		t.Errorf("Number of transactions returned is invalid")
	}

	for _, transaction := range transactions {
		if transaction.Amount != 1 {
			t.Errorf("Transaction amount is invalid")
		}
	}
}

func TestGetTransactionsForNonExistentAccount(t *testing.T) {
	transactions, err := TransactionDB.Get(12345678)
	assertNoError(t, err)
	if len(transactions) != 0 {
		t.Errorf("Transactions were returned and should not have been")
	}
}

func TestSetValidTransaction(t *testing.T) {
	now := time.Now().Unix()
	transaction := atm.Transaction{
		AccountID: 987654321,
		DateTime:  now,
		Amount:    1.00,
		Balance:   2.00,
	}
	err := TransactionDB.Set(transaction)
	assertNoError(t, err)

	transactions, err := TransactionDB.Get(transaction.AccountID)
	assertNoError(t, err)

	// Validate transaction now exists in the transaction database using the time as the unique ID
	transactionFound := false
	for _, transaction := range transactions {
		if transaction.DateTime == now {
			transactionFound = true
		}
	}

	if !transactionFound {
		t.Errorf("Transaction was created but could not be found in database")
	}
}

func TestGetTransactionsInvalidDatabase(t *testing.T) {
	_, err := InvalidTransactionDB.Get(12345678)
	assertError(t, err)
}
