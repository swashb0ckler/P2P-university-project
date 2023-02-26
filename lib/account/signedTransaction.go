package account

import (
	"fmt"
	"math/big"
	"p2p/lib/encrypt"
	"strconv"
)

type SignedTransaction struct {
	//ID        string // Any string
	From      string // A verification key coded as a string To string // A verification key coded as a string Amount int // Amount to transfer
	To        string
	Amount    int
	Signature string // Potential signature coded as string
}

func (l *Ledger) SignedTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// -- Converting variables to big Ints for verifySignature() --

	// signature to big int
	signatureBigInt := new(big.Int)
	signatureBigInt, error := signatureBigInt.SetString(t.Signature, 10)
	if !error {
		fmt.Println("SetString: error")
		return
	}

	// modulus to big int
	modulus := new(big.Int)
	modulus, error2 := modulus.SetString(t.From, 10)
	if !error2 {
		fmt.Println("SetString: error2")
		return
	}

	// message collected in an allFields message and then converted to big int
	amountString := strconv.Itoa(t.Amount)
	allFieldsString := t.From + t.To + amountString
	message := new(big.Int)
	message, error3 := message.SetString(allFieldsString, 10)
	if !error3 {
		fmt.Println("SetString: error3")
		return
	}

	// Making e
	e := big.NewInt(3)

	// Verify signature
	validSignature := encrypt.VerifySignature(message, e, signatureBigInt, modulus)

	if validSignature {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount
	}
}
