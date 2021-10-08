package atm_test

import (
	"testing"
	"time"

	"github.com/AndrewCopeland/atm"
)

func TestSession(t *testing.T) {
	session := &atm.Session{}

	// Authorize then validate session is active then logout of session and validate session is not inasctive
	session.Authorize(defaultAccount.AccountID)
	err := session.Valid(defaultAccount.AccountID)
	assertNoError(t, err)
	err = session.Valid(98765)
	assertErrorIsError(t, err, atm.ErrSessionInvalidAccountID)
	err = session.LogOut()
	assertNoError(t, err)
	err = session.Valid(defaultAccount.AccountID)
	assertErrorIsError(t, err, atm.ErrSessionNoActiveSession)

	// Test session is inactive after 2 mins
	session.Authorize(defaultAccount.AccountID)
	session.LastActivity = time.Now().Unix() - 121
	err = session.Valid(defaultAccount.AccountID)
	assertErrorIsError(t, err, atm.ErrSessionTimedOut)

	// Logout of the timed out session
	err = session.LogOut()
	assertErrorIsError(t, err, atm.ErrSessionTimedOut)
}
