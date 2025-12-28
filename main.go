package main

import "fmt"

var overdraftLimit = -500.0

type BankAccount struct {
	balance float64
}

func (account *BankAccount) Withdraw(amount float64) bool {
	if account.balance-amount < overdraftLimit {
		return false
	}
	account.balance -= amount
	return true
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
	account   *BankAccount
	action    Action
	amount    float64
	succeeded bool
}

func (c *BankAccountCommand) Call() {
	switch c.action {
	case Deposit:
		c.account.Deposit(c.amount)
		c.succeeded = true
	case Withdraw:
		c.succeeded = c.account.Withdraw(c.amount)
	}
}

func (c *BankAccountCommand) Undo() {
	if !c.succeeded {
		return
	}
	switch c.action {
	case Deposit:
		c.account.Withdraw(c.amount)
	case Withdraw:
		c.account.Deposit(c.amount)
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
	cmd.Undo()
	fmt.Println("Account balance after undoing withdrawal:", account.balance)
	cmd2.Undo()
	fmt.Println("Account balance after undoing deposit:", account.balance)
}
