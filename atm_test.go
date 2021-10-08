package atm_test

import (
	"errors"
	"testing"
	"time"

	"github.com/AndrewCopeland/atm"
)

type AccountDBTest struct {
	getAccount      atm.Account
	getAccountError error
	setAccountError error
}

func (a AccountDBTest) Get(accountID int) (atm.Account, error) {
	if accountID != a.getAccount.AccountID {
		return atm.Account{}, atm.ErrAccountNotFound
	}
	return a.getAccount, a.getAccountError
}

func (a AccountDBTest) Set(account atm.Account) error {
	return a.setAccountError
}

type TransactionDBTest struct {
	getTransactions      []atm.Transaction
	getTransactionsError error
	setTransactionError  error
}

func (t TransactionDBTest) Get(accountID int) ([]atm.Transaction, error) {
	return t.getTransactions, t.getTransactionsError
}

func (t TransactionDBTest) Set(transaction atm.Transaction) error {
	return t.setTransactionError
}

var defaultAccount = atm.Account{
	AccountID: 12345678,
	PIN:       "1234",
	Balance:   100.00,
}

var defaultAccountDB = AccountDBTest{
	getAccount:      defaultAccount,
	getAccountError: nil,
}

var defaultTranscationDB = TransactionDBTest{
	getTransactions: []atm.Transaction{
		{
			AccountID: 12345678,
			DateTime:  time.Now().Unix(),
			Amount:    1.00,
			Balance:   1.00,
		},
	},
}

func newTestATM(accountDB AccountDBTest, transactionDB TransactionDBTest) atm.ATM {
	atm := atm.ATM{
		AccountDB:     accountDB,
		TransactionDB: transactionDB,
		ATMBalance:    200.00,
		Session:       &atm.Session{},
	}

	return atm
}

func defaultTestATM() atm.ATM {
	return newTestATM(defaultAccountDB, defaultTranscationDB)
}

func TestAuthorize(t *testing.T) {
	atm := defaultTestATM()

	// Test invalid PIN
	result := atm.Authorize(defaultAccount.AccountID, "0000")
	if result {
		t.Error("Authenticated but should not have authenticated")
	}

	// Test invalid database
	invalidDatabase := defaultAccountDB
	invalidDatabase.getAccountError = errors.New("Failed to get account")
	atm = newTestATM(invalidDatabase, defaultTranscationDB)
	result = atm.Authorize(defaultAccount.AccountID, defaultAccount.PIN)
	if result {
		t.Error("Authenticated successful but database is not working")
	}
	atm = defaultTestATM()

	// Test valid PIN
	result = atm.Authorize(defaultAccount.AccountID, defaultAccount.PIN)
	if !result {
		t.Error("Failure to authorize")
	}
	err := atm.Session.Valid(defaultAccount.AccountID)
	assertNoError(t, err)
}

func TestWithdraw(t *testing.T) {
	testATM := defaultTestATM()

	// no active session
	_, err := testATM.Withdraw(defaultAccount.AccountID, 20)
	assertErrorIsError(t, err, atm.ErrSessionNoActiveSession)

	// activate session
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}

	// amount not divisable by 20
	_, err = testATM.Withdraw(defaultAccount.AccountID, 1)
	assertErrorIsError(t, err, atm.ErrWithdrawAmountNoMultipleOf20)

	// atm has no funds
	testATM.ATMBalance = 0
	_, err = testATM.Withdraw(defaultAccount.AccountID, 20)
	assertErrorIsError(t, err, atm.ErrWithdrawATMNoFunds)
	testATM.ATMBalance = 100.00

	// atm has insufficent funds
	_, err = testATM.Withdraw(defaultAccount.AccountID, 120)
	assertErrorIsError(t, err, atm.ErrWithdrawATMInsufficientFunds)

	// account is overdrawn
	accountDB := AccountDBTest{
		getAccount: atm.Account{
			AccountID: defaultAccount.AccountID,
			PIN:       "1234",
			Balance:   -5.0,
		},
		getAccountError: nil,
	}
	overDrawnATM := newTestATM(accountDB, defaultTranscationDB)
	overDrawnATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	_, err = overDrawnATM.Withdraw(defaultAccount.AccountID, 100)
	assertErrorIsError(t, err, atm.ErrWithdrawAccountOverdrawn)

	// overdrawn but withdraw
	testATM.ATMBalance = 200.00
	overdrawn, err := testATM.Withdraw(defaultAccount.AccountID, 120)
	assertNoError(t, err)
	if !overdrawn {
		t.Errorf("Expected overdrawn")
	}

	// valid withdraw
	overdrawn, err = testATM.Withdraw(defaultAccount.AccountID, 20)
	if overdrawn {
		t.Errorf("Unexpected overdrawn")
	}
	assertNoError(t, err)
}

func TestDepositValid(t *testing.T) {
	testATM := defaultTestATM()

	// Test deposit no session
	err := testATM.Deposit(defaultAccount.AccountID, 1)
	assertErrorIsError(t, err, atm.ErrSessionNoActiveSession)

	// Test deposit valid session
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	err = testATM.Deposit(defaultAccount.AccountID, 1)
	assertNoError(t, err)
}

func TestBalanceValid(t *testing.T) {
	testATM := defaultTestATM()

	// test balance no session
	_, err := testATM.Balance(defaultAccount.AccountID)
	assertErrorIsError(t, err, atm.ErrSessionNoActiveSession)

	// test balance with active session
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	balance, err := testATM.Balance(defaultAccount.AccountID)
	if err != nil {
		t.Error("Failed to read balance")
	}

	if balance != defaultAccount.Balance {
		t.Error("Balance is not correct")
	}
}

func TestHistoryValid(t *testing.T) {
	testATM := defaultTestATM()

	// Get history with no session
	transactions := testATM.History(defaultAccount.AccountID)
	if len(transactions) != 0 {
		t.Errorf("Retrieved history of account even though no active sessions")
	}

	// Get history with active session
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	transactions = testATM.History(defaultAccount.AccountID)
	if len(transactions) < 1 {
		t.Error("Failed to get history of account")
	}
}

func TestLogout(t *testing.T) {
	testATM := defaultTestATM()

	// logout of atm without active session
	err := testATM.Logout()
	assertErrorIsError(t, err, atm.ErrSessionNoActiveSession)

	// logout of atm with active session
	authorized := testATM.Authorize(defaultAccount.AccountID, defaultAccount.PIN)
	if !authorized {
		t.Error("Failed to authorize to ATM")
	}
	err = testATM.Logout()
	assertNoError(t, err)
}
