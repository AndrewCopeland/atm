package atm

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CommandRun func(*ATM, []string) error

type ICommand interface {
	Name() string
	Run(*ATM, []string) error
}

type Command struct {
	name     string
	usage    string
	function CommandRun
}

// Name shows the name of the command
func (c Command) Name() string {
	return c.name
}

// Run executes the command with the provided arguments
// An error is returned when the command failed to run
func (c Command) Run(atm *ATM, args []string) error {
	err := c.function(atm, args)

	// Append usage if invalid command was provided
	if err == ErrConsoleInvalidCommand {
		return errors.New(ErrConsoleInvalidCommand.Error() + c.usage)
	}
	return err
}

// DefaultCommands returns all of the commands valid in the console
func DefaultCommands() []Command {
	return []Command{
		{
			name:  "authorize",
			usage: "authorize <account_id> <pin>",
			function: func(atm *ATM, args []string) error {
				if len(args) != 3 {
					return ErrConsoleInvalidCommand
				}

				accountID, err := strconv.Atoi(args[1])
				if err != nil {
					return ErrAccountIDNotInteger
				}

				pin := args[2]

				authorized := atm.Authorize(accountID, pin)
				if !authorized {
					return ErrConsoleAuthorizationFailed
				}
				fmt.Printf("%d successfully authorized.\n", accountID)
				return nil
			},
		},
		{
			name:  "withdraw",
			usage: "withdraw <amount>",
			function: func(atm *ATM, args []string) error {
				if len(args) != 2 {
					return ErrConsoleInvalidCommand
				}

				amount, err := strconv.Atoi(args[1])
				if err != nil {
					return ErrConsoleInvalidAmount
				}

				overdrawn, err := atm.Withdraw(atm.Session.AccountID, amount)
				if err != nil {
					return err
				}

				newBalance, err := atm.Balance(atm.Session.AccountID)
				if err != nil {
					return err
				}

				fmt.Printf("Amount dispensed: %d\n", amount)
				if overdrawn {
					fmt.Printf("You have been charged an overdraft fee of $5. Current balance:  %.2f\n", newBalance)
				} else {
					fmt.Printf("Current balance: %.2f\n", newBalance)
				}
				return nil
			},
		},
		{
			name:  "deposit",
			usage: "deposit <amount>",
			function: func(atm *ATM, args []string) error {
				if len(args) != 2 {
					return ErrConsoleInvalidCommand
				}

				amount, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					return ErrConsoleInvalidAmount
				}

				err = atm.Deposit(atm.Session.AccountID, amount)
				if err != nil {
					return err
				}

				balance, err := atm.Balance(atm.Session.AccountID)
				if err != nil {
					return err
				}

				fmt.Printf("Current balance: %.2f\n", balance)
				return nil
			},
		},
		{
			name:  "balance",
			usage: "balance",
			function: func(atm *ATM, args []string) error {
				balance, err := atm.Balance(atm.Session.AccountID)
				if err != nil {
					return err
				}

				fmt.Printf("Current balance: %.2f\n", balance)
				return nil
			},
		},
		{
			name:  "history",
			usage: "history",
			function: func(atm *ATM, args []string) error {
				transactions := atm.History(atm.Session.AccountID)
				if len(transactions) == 0 {
					fmt.Println("No history found")
					return nil
				}

				for i := len(transactions) - 1; i >= 0; i-- {
					t := time.Unix(transactions[i].DateTime, 0)
					fmt.Printf("%s %.2f %.2f\n", t.Format("01-02-2006 15:04:05"), transactions[i].Amount, transactions[i].Balance)
				}
				return nil
			},
		},
		{
			name:  "logout",
			usage: "logout",
			function: func(atm *ATM, args []string) error {
				accountID := atm.Session.AccountID
				err := atm.Logout()
				if err != nil {
					return errors.New("No account is currently authorized.")
				}
				fmt.Printf("Account %d logged out.\n", accountID)
				return nil
			},
		},
		{
			name:  "end",
			usage: "end",
			function: func(atm *ATM, args []string) error {
				os.Exit(0)
				return nil
			},
		},
	}
}

// RunCommand will execute the command given the command string
// If an invalid command is used then the help message is displayed
func RunCommand(atm *ATM, command string) error {
	args := strings.Split(command, " ")
	if len(args) == 0 || args[0] == "" {
		return nil
	}

	commands := DefaultCommands()
	found := false
	for _, c := range commands {
		if c.Name() == strings.ToLower(args[0]) {
			found = true
			return c.Run(atm, args)
		}
	}

	if !found {
		fmt.Println("Command Usage:")
		for _, c := range commands {
			fmt.Println(c.usage)
		}
		return errors.New("Invalid command")
	}

	return nil
}
