package atm

import (
	"time"
)

type ISession interface {
	Authorize(int)
	TimedOut() bool
	Refresh()
	Logout() bool
	Valid(int) error
}

type Session struct {
	// When the session started in epoch
	LastActivity int64

	AccountID int
}

// Authorize will set the LastActivity time to now and the AccountID of the session
func (s *Session) Authorize(accountID int) {
	s.AccountID = accountID
	s.Refresh()
}

// Refresh updates the LastActivity to now
func (s *Session) Refresh() {
	s.LastActivity = time.Now().Unix()
}

// TimedOut checks if the session has timed out after 2 mins
func (s *Session) TimedOut() bool {
	difference := time.Now().Unix() - s.LastActivity
	// if activity has not happened in 2 mins or more
	if difference > 120 {
		return true
	}
	return false
}

// LogOut will logout of the session
func (s *Session) LogOut() error {
	err := s.Valid(s.AccountID)
	if err != nil {
		return err
	}

	s.AccountID = 0
	s.LastActivity = 0
	return nil
}

// Valid will validate the session exists and not timed out and will refresh the LastActivity
func (s *Session) Valid(accountID int) error {
	if s.AccountID == 0 || s.LastActivity == 0 {
		return ErrSessionNoActiveSession
	}
	if s.AccountID != accountID {
		return ErrSessionInvalidAccountID
	}
	if s.TimedOut() {
		return ErrSessionTimedOut
	}

	s.Refresh()
	return nil
}
