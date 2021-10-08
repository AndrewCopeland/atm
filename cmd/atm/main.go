package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AndrewCopeland/atm"
)

func main() {
	accountDB := atm.AccountDB{
		DBFile: "./accounts.csv",
	}
	transactionDB := atm.TransactionDB{
		DBFile: "./transactions.csv",
	}
	a := &atm.ATM{
		AccountDB:     accountDB,
		TransactionDB: transactionDB,
		ATMBalance:    10000.00,
		Session:       &atm.Session{},
	}

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := atm.RunCommand(a, scanner.Text())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("> ")
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}
}
