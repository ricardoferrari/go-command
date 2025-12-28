package main

import "fmt"

var overdraftLimit = -500.0

type BankAccount struct {
	balance float64
}

func (account *BankAccount) Withdraw(amount float64) {
	if account.balance-amount < overdraftLimit {
		return
	}
	account.balance -= amount
}

func (account *BankAccount) Deposit(amount float64) {
	account.balance += amount
}

type Command interface {
	Call()
	Undo()
}

type Action int

const (
	Deposit Action = iota
	Withdraw
)

type BankAccountCommand struct {
	account *BankAccount
	action  Action
	amount  float64
}

func (c *BankAccountCommand) Call() {
	switch c.action {
	case Deposit:
		c.account.Deposit(c.amount)
	case Withdraw:
		c.account.Withdraw(c.amount)
	}
}

func main() {
	account := &BankAccount{balance: 1000}
	cmd := &BankAccountCommand{account: account, action: Withdraw, amount: 200}
	cmd.Call()
	fmt.Println("Account balance:", account.balance)
	cmd2 := &BankAccountCommand{account: account, action: Deposit, amount: 500}
	cmd2.Call()
	fmt.Println("Account balance after deposit:", account.balance)
}
