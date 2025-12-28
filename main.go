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
	Succeeded() bool
	SetSucceeded(value bool)
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

func (c *BankAccountCommand) Succeeded() bool {
	return c.succeeded
}

func (c *BankAccountCommand) SetSucceeded(value bool) {
	c.succeeded = value
}

type CompositeBankAccountCommand struct {
	commands []Command
}

func (c *CompositeBankAccountCommand) Call() {
	for _, cmd := range c.commands {
		cmd.Call()
	}
}

func (c *CompositeBankAccountCommand) Undo() {
	for i := len(c.commands) - 1; i >= 0; i-- {
		c.commands[i].Undo()
	}
}

func (c *CompositeBankAccountCommand) Succeeded() bool {
	for _, cmd := range c.commands {
		if !cmd.Succeeded() {
			return false
		}
	}
	return true
}

func (c *CompositeBankAccountCommand) SetSucceeded(value bool) {
	for _, cmd := range c.commands {
		cmd.SetSucceeded(value)
	}
}

type MoneyTransferCommand struct {
	CompositeBankAccountCommand
	from   *BankAccount
	to     *BankAccount
	amount float64
}

func NewMoneyTransferCommand(from, to *BankAccount, amount float64) *MoneyTransferCommand {
	withdrawCmd := &BankAccountCommand{account: from, action: Withdraw, amount: amount}
	depositCmd := &BankAccountCommand{account: to, action: Deposit, amount: amount}
	commands := []Command{withdrawCmd, depositCmd}
	return &MoneyTransferCommand{
		CompositeBankAccountCommand: CompositeBankAccountCommand{commands: commands},
		from:                        from,
		to:                          to,
		amount:                      amount,
	}
}

func (c *MoneyTransferCommand) Call() {
	succeded := true
	for _, cmd := range c.commands {
		if succeded {
			cmd.Call()
			succeded = cmd.Succeeded()
		} else {
			cmd.SetSucceeded(false)
		}
	}
}

func main() {
	// Simple bank account command example
	fmt.Println("Simple Bank Account Command Example:")
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

	// Money transfer example
	fmt.Println("\nMoney Transfer Command Example:")
	accountA := &BankAccount{balance: 1000}
	accountB := &BankAccount{balance: 500}
	transferCmd := NewMoneyTransferCommand(accountA, accountB, 300)
	transferCmd.Call()
	fmt.Println("Account A balance after transfer:", accountA.balance)
	fmt.Println("Account B balance after transfer:", accountB.balance)
	// Print whether the transfer succeeded
	fmt.Println("Did the transfer succeed?", transferCmd.Succeeded())
	transferCmd.Undo()
	fmt.Println("Account A balance after undoing transfer:", accountA.balance)
	fmt.Println("Account B balance after undoing transfer:", accountB.balance)

	// Composite command example exceeding overdraft limit
	fmt.Println("\nComposite Command Exceeding Overdraft Limit Example:")
	largeTransferCmd := NewMoneyTransferCommand(accountA, accountB, 2000)
	largeTransferCmd.Call()
	fmt.Println("Account A balance after large transfer attempt:", accountA.balance)
	fmt.Println("Account B balance after large transfer attempt:", accountB.balance)
	// Print whether the large transfer succeeded
	fmt.Println("Did the large transfer succeed?", largeTransferCmd.Succeeded())
	largeTransferCmd.Undo()
	fmt.Println("Account A balance after undoing large transfer attempt:", accountA.balance)
	fmt.Println("Account B balance after undoing large transfer attempt:", accountB.balance)
}
