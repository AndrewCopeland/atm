package atm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AndrewCopeland/atm"
)

func TestConsoleAuthorize(t *testing.T) {
	a := defaultTestATM()
	testATM := &a

	// Invalid authorize command
	err := atm.RunCommand(testATM, "authorize")
	assertErrorContains(t, err, atm.ErrConsoleInvalidCommand.Error())

	// authorize with invalid account ID
	err = atm.RunCommand(testATM, "authorize string 6783")
	assertErrorIsError(t, err, atm.ErrAccountIDNotInteger)

	// authorize successfully
	err = atm.RunCommand(testATM, fmt.Sprintf("authorize %d %s", defaultAccount.AccountID, defaultAccount.PIN))
	assertNoError(t, err)
}

func TestConsoleWithdraw(t *testing.T) {
	a := defaultTestATM()
	testATM := &a

	// invalid withdraw command
	err := atm.RunCommand(testATM, "withdraw")
	assertErrorContains(t, err, atm.ErrConsoleInvalidCommand.Error())

	// withdraw with invalid amount
	err = atm.RunCommand(testATM, "withdraw abc")
	assertErrorIsError(t, err, atm.ErrConsoleInvalidAmount)

	// with with valid amount
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	err = atm.RunCommand(testATM, "withdraw 20")
	assertNoError(t, err)
}

func TestConsoleDeposit(t *testing.T) {
	a := defaultTestATM()
	testATM := &a

	// invalid deposit command
	err := atm.RunCommand(testATM, "deposit")
	assertErrorContains(t, err, atm.ErrConsoleInvalidCommand.Error())

	// invalid deposit amount
	err = atm.RunCommand(testATM, "deposit invalid")
	assertErrorIsError(t, err, atm.ErrConsoleInvalidAmount)

	// successful deposit
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}
	err = atm.RunCommand(testATM, "deposit 20")
	assertNoError(t, err)
}

func TestConsoleBalance(t *testing.T) {
	a := defaultTestATM()
	testATM := &a
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}

	err := atm.RunCommand(testATM, "balance")
	assertNoError(t, err)
}

func TestConsoleHistory(t *testing.T) {
	a := defaultTestATM()
	testATM := &a
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}

	err := atm.RunCommand(testATM, "history")
	assertNoError(t, err)
}

func TestConsoleLogoff(t *testing.T) {
	a := defaultTestATM()
	testATM := &a
	testATM.Session = &atm.Session{
		LastActivity: time.Now().Unix(),
		AccountID:    defaultAccount.AccountID,
	}

	err := atm.RunCommand(testATM, "logout")
	assertNoError(t, err)
}
