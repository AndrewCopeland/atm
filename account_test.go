package atm_test

import (
	"testing"

	"github.com/AndrewCopeland/atm"
)

var AccountDB = atm.AccountDB{
	DBFile: "./accounts_test.csv",
}

var InvalidAccountDB = atm.AccountDB{
	DBFile: "./accounts_not_real.csv",
}

func TestGetValidAccount(t *testing.T) {
	account, err := AccountDB.Get(7089382418)
	assertNoError(t, err)

	// 7089382418,0075,0.00
	if account.AccountID != 7089382418 {
		t.Errorf("Account ID is incorrect")
	}
	if account.PIN != "0075" {
		t.Errorf("Account PIN is incorrect")
	}
	if account.Balance != 0.00 {
		t.Errorf("Account balance is incorrect")
	}
}

func TestGetInvalidAccount(t *testing.T) {
	_, err := AccountDB.Get(123456789)
	assertError(t, err)
}

func TestSetValidAccount(t *testing.T) {
	// 2859459814,7386,10.24
	account := atm.Account{
		AccountID: 2859459814,
		PIN:       "7386",
		Balance:   5.24,
	}
	err := AccountDB.Set(account)
	assertNoError(t, err)

	// Validate that the account was updated by retrieving it again
	resultAccount, err := AccountDB.Get(account.AccountID)
	assertNoError(t, err)

	if account.AccountID != resultAccount.AccountID {
		t.Errorf("Account ID is incorrect")
	}
	if account.PIN != resultAccount.PIN {
		t.Errorf("Account PIN is incorrect")
	}
	if account.Balance != resultAccount.Balance {
		t.Errorf("Account balance is incorrect")
	}
}

func TestSetInvalidAccount(t *testing.T) {
	account := atm.Account{
		AccountID: 123456789,
		PIN:       "1234",
		Balance:   10.21,
	}
	err := AccountDB.Set(account)
	assertError(t, err)
}

func TestGetAccountInvalidDatabaseFile(t *testing.T) {
	_, err := InvalidAccountDB.Get(12345543)
	assertError(t, err)
}
