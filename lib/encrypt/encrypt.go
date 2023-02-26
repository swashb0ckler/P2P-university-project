package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
)

// Outputted keys

func Encrypt(message, e, n *big.Int) *big.Int {
	exponent := big.NewInt(0).Exp(message, e, n)
	return exponent
}

func Decrypt(cipherText, d, n *big.Int) *big.Int {
	exponent2 := big.NewInt(0).Exp(cipherText, d, n) //n instead of nil
	return exponent2
}

// parameter k needs to be bigger than 10
func GeneratePQ(k int) (*big.Int, *big.Int, *big.Int) {

	p := GeneratePrime(k)
	q := GeneratePrime(k)

	pMinus1 := big.NewInt(0).Sub(p, big.NewInt(1))
	qMinus1 := big.NewInt(0).Sub(q, big.NewInt(1))

	n := big.NewInt(0).Mul(p, q)

	//Converts n to a bitstring
	bitString := fmt.Sprintf("%b", n)

	//Check if k is equal to length of n
	if k == len(bitString) {
		fmt.Println("n is the same length as the value of k")
	}

	// Carmichael function

	x := big.NewInt(0)
	x.GCD(nil, nil, pMinus1, qMinus1)

	//Calculating d, the private key
	e := big.NewInt(3)
	d := big.NewInt(0).ModInverse(e, big.NewInt(0).Mul(pMinus1, qMinus1))

	return n, e, d

}

// Generates a random prime which fullfils GCD(e, p-1) = GCD(e, q-1) = 1
func GeneratePrime(a int) *big.Int {
	e := big.NewInt(3)
	randNumber, _ := rand.Prime(rand.Reader, a/2)
	number1 := big.NewInt(1)
	tempVar := big.NewInt(0)
	gcdMinus1 := tempVar.GCD(nil, nil, e, big.NewInt(0).Sub(randNumber, big.NewInt(1)))
	// while GCD of e and randNumber != 1 generate new random prime
	for gcdMinus1.Cmp(number1) != 0 {
		randNumber, _ = rand.Prime(rand.Reader, a/2)
		gcdMinus1 = tempVar.GCD(nil, nil, e, big.NewInt(0).Sub(randNumber, big.NewInt(1)))

	}
	return randNumber
}

func CreateSignature(message, d, n *big.Int) *big.Int {

	// Saves hashValue in a variable
	hashValue := HashingFunc(message)

	//Takes the hashvalue^d and saves in a variable
	result := big.NewInt(0).Exp(hashValue, d, n)

	// Measuring the time of RSA signature ---

	//Returns the signature s
	return result

}

func VerifySignature(message, e, s, n *big.Int) bool {

	// Saves s^e mod n in a variable
	result := big.NewInt(0).Exp(s, e, n)

	//Checks if the
	if HashingFunc(message).Cmp(result) == 0 { // n
		return true
	} else {
		return false
	}
}

func HashingFunc(message *big.Int) *big.Int {

	//Turns message into bytes
	var b = message.Bytes()

	//Creates sha object and writes the messages bytes to it
	shaObj := sha256.New()
	shaObj.Write(b)

	//Initializes new big int and takes the check sum of the msg bytes and turns it into a big int
	hashValue := big.NewInt(0)
	hashValue = hashValue.SetBytes(shaObj.Sum(nil))

	return hashValue

}

func EncryptToFile(fileName string, m *big.Int, cipherKey []byte) {
	plaintext := m.Bytes()
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	IV := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, IV); err != nil {
		panic(err)
	}

	CTRstream := cipher.NewCTR(block, IV)
	CTRstream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	ioutil.WriteFile(fileName, []byte(ciphertext), 0777)
}

func DecryptFromFile(fileName string, cipherKey []byte) []byte {
	cipherstring, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	ciphertext := []byte(cipherstring)

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		panic(err)
	}

	IV := ciphertext[:aes.BlockSize]
	CTRstream := cipher.NewCTR(block, IV)
	plainText := make([]byte, len(ciphertext[aes.BlockSize:]))
	CTRstream.XORKeyStream(plainText, ciphertext[aes.BlockSize:])
	return plainText
}
