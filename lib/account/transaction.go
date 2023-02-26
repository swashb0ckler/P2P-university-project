package account

import "fmt"

type Transaction struct {
	//ID     string
	From   string
	To     string
	Amount int
}

func (l *Ledger) Transaction(t *Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Println("sending amount", t.Amount, "from", t.From, "to", t.To)
	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
}
