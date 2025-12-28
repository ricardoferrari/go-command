# Command Pattern in Go

This project demonstrates the **Command Pattern**, a behavioral design pattern that encapsulates a request as an object, allowing you to parameterize clients with different requests, queue or log requests, and support undoable operations.

## Overview

The Command Pattern implementation in this project uses a banking system to demonstrate three key aspects of the pattern:
1. **Simple Command** - Basic command execution
2. **Undo Functionality** - Reversible operations
3. **Composite Commands** - Complex operations composed of multiple commands

---

## 1. Simple Command

### Concept
A simple command encapsulates a single action that can be executed independently. The pattern decouples the object that invokes the operation from the object that performs it.

### Implementation

The `Command` interface defines the contract for all commands:

```go
type Command interface {
    Call()              // Execute the command
    Undo()              // Reverse the command
    Succeeded() bool    // Check if command succeeded
    SetSucceeded(value bool)
}
```

The `BankAccountCommand` implements this interface for basic banking operations:

```go
type BankAccountCommand struct {
    account   *BankAccount
    action    Action      // Deposit or Withdraw
    amount    float64
    succeeded bool
}
```

### Usage Example

```go
account := &BankAccount{balance: 1000}
cmd := &BankAccountCommand{account: account, action: Withdraw, amount: 200}
cmd.Call()  // Executes the withdrawal
```

### Key Benefits
- **Encapsulation**: The command encapsulates all information needed to perform the action
- **Parameterization**: Different actions (Deposit/Withdraw) use the same structure
- **Deferred Execution**: Commands can be created and executed later

---

## 2. Undo and Redo

### Concept
One of the most powerful features of the Command Pattern is the ability to reverse operations. Each command knows how to undo its own action.

### Implementation

The `Undo()` method reverses the operation performed by `Call()`:

```go
func (c *BankAccountCommand) Undo() {
    if !c.succeeded {
        return  // Don't undo if command didn't succeed
    }
    switch c.action {
    case Deposit:
        c.account.Withdraw(c.amount)  // Reverse deposit with withdrawal
    case Withdraw:
        c.account.Deposit(c.amount)   // Reverse withdrawal with deposit
    }
}
```

### Usage Example

```go
cmd := &BankAccountCommand{account: account, action: Withdraw, amount: 200}
cmd.Call()   // Balance: 800
cmd.Undo()   // Balance: 1000 (restored)
```

### Key Features
- **State Tracking**: The `succeeded` flag ensures only successful operations are undone
- **Symmetry**: Each action has a clear inverse operation
- **Safety**: Failed operations are not undone to maintain consistency

---

## 3. Composite Commands

### Concept
Composite commands combine multiple commands into a single operation. This allows complex transactions to be treated as a single unit while maintaining the ability to undo the entire operation.

### Implementation

The `CompositeBankAccountCommand` groups multiple commands:

```go
type CompositeBankAccountCommand struct {
    commands []Command
}
```

The `MoneyTransferCommand` extends the composite command to implement atomic money transfers:

```go
type MoneyTransferCommand struct {
    CompositeBankAccountCommand
    from   *BankAccount
    to     *BankAccount
    amount float64
}
```

A money transfer consists of two operations:
1. **Withdraw** from the source account
2. **Deposit** to the destination account

### Key Implementation Details

**Sequential Execution with Failure Handling**:
```go
func (c *MoneyTransferCommand) Call() {
    succeded := true
    for _, cmd := range c.commands {
        if succeded {
            cmd.Call()
            succeded = cmd.Succeeded()
        } else {
            cmd.SetSucceeded(false)  // Mark subsequent commands as failed
        }
    }
}
```

**Reverse-Order Undo**:
```go
func (c *CompositeBankAccountCommand) Undo() {
    for i := len(c.commands) - 1; i >= 0; i-- {
        c.commands[i].Undo()  // Undo in reverse order
    }
}
```

### Usage Example

```go
accountA := &BankAccount{balance: 1000}
accountB := &BankAccount{balance: 500}
transferCmd := NewMoneyTransferCommand(accountA, accountB, 300)
transferCmd.Call()  // Transfer 300 from A to B
// A: 700, B: 800

transferCmd.Undo()  // Reverse the entire transfer
// A: 1000, B: 500 (restored)
```

### Key Benefits
- **Atomicity**: The entire operation succeeds or fails as a unit
- **Transaction Safety**: If withdrawal fails (e.g., overdraft limit), deposit won't execute
- **Complete Rollback**: Undo reverses all sub-commands in the correct order
- **Extensibility**: New composite operations can be built from existing commands

### Failure Handling Example

When a transfer exceeds the overdraft limit:

```go
largeTransferCmd := NewMoneyTransferCommand(accountA, accountB, 2000)
largeTransferCmd.Call()
fmt.Println(largeTransferCmd.Succeeded())  // false

// Neither account is modified because withdrawal failed
// Deposit is marked as failed and not executed
```

---

## Running the Project

Execute the demonstration:

```bash
go run main.go
```

The output shows:
1. Simple commands (deposit/withdraw with undo)
2. Successful money transfers with undo
3. Failed transfers respecting overdraft limits

---

## Design Pattern Advantages

1. **Separation of Concerns**: The invoker doesn't need to know how commands are executed
2. **Undo/Redo Support**: Commands maintain enough state to reverse themselves
3. **Command Queuing**: Commands can be queued, logged, or scheduled
4. **Macro Commands**: Complex operations can be built from simpler ones
5. **Transaction Support**: Multiple operations can be treated as atomic units

## Further Extensions

This pattern can be extended to support:
- **Command History**: Store executed commands for multi-level undo/redo
- **Command Queuing**: Execute commands asynchronously
- **Command Logging**: Persist commands for auditing or replay
- **Transactional Systems**: Ensure atomicity in distributed systems
