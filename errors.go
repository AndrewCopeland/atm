package atm

import (
	"errors"
)

// authorize errors
var (
	ErrAuthorizationUnsuccessful = errors.New("Authorization failed.")
	ErrAuthorizationRequired     = errors.New("Authorization required.")
)

// withdraw errors
var (
	ErrWithdrawATMInsufficientFunds = errors.New("Unable to dispense full amount requested at this time.")
	ErrWithdrawATMNoFunds           = errors.New("Unable to process your withdrawal at this time.")
	ErrWithdrawAccountOverdrawn     = errors.New("Your account is overdrawn! You may not make withdrawals at this time.")
	ErrWithdrawAmountNoMultipleOf20 = errors.New("Unable to process since amount is not a multiple of 20.")
)

// logout errors
var (
	ErrLogoutNoActiveSession = errors.New("No account is currently authorized.")
)

// account db error
var (
	ErrAccountNotFound        = errors.New("Account could not be found in database.")
	ErrAccountIDNotInteger    = errors.New("Account ID is not an integer")
	ErrAccountBalanceNotFloat = errors.New("Balance is not a float")
)

// transaction db error
var (
	ErrTransactionAccountIDNotInteger = errors.New("Account ID is not an integer")
	ErrTransactionDateTimeNotInteger  = errors.New("Datetime is not an integer")
	ErrTransactionAmounteNotFloat     = errors.New("Amount is not a float")
	ErrTransactionBalanceNotFloat     = errors.New("Balance is not a float")
)

// session error
var (
	ErrSessionNoActiveSession  = errors.New("No active session found. Authorization required.")
	ErrSessionInvalidAccountID = errors.New("Invalid account ID for session.")
	ErrSessionTimedOut         = errors.New("Session has timed out.")
)

// console error
var (
	ErrConsoleInvalidCommand      = errors.New("Invalid command. e.g. ")
	ErrConsoleAuthorizationFailed = errors.New("Authorization failed.")
	ErrConsoleInvalidAmount       = errors.New("Amount is not a valid")
)
