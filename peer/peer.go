package peer

import (
	"encoding/gob"
	"fmt"
	"math/big"
	"net"
	"p2p/lib/account"
	"p2p/lib/encrypt"
	"p2p/lib/util"
	"strconv"
)

type Peer struct {
	ID                      string
	Ip                      string
	Port                    string
	Peerlist                PeerList
	ConnectionEncoderStruct ConnectionEncoderStruct //ConnectionEncoders []*gob.Encoder
	Ledger                  account.Ledger
	SecretKey               *big.Int
}

type PeerList struct {
	Peers []string
}

type ConnectionEncoderStruct struct {
	ConnectionEncoderList []*gob.Encoder
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// Create a new Peer
func InitPeer(address string) (*Peer, string, string) {
	// 0. Create a Peer object
	p := &Peer{}

	// 1. Start listening
	ip, port := p.listen()
	p.Ip = ip
	p.Port = port

	// key generation
	n, _, d := encrypt.GeneratePQ(1000)

	// Setting ID to public key and inserting d into peer struct
	p.SecretKey = d
	p.ID = n.String() // p.Port before

	// 2. Create ledger
	p.Ledger = *account.MakeLedger()
	p.Ledger.Accounts[p.ID] = 0

	// 3. Attempt to connect
	if address != "" {
		p.connect(address, true, true)
	} else {
		//do nothing
	}

	// 4. Return the Peer object and the IP and port that it is listening on
	return p, ip, port
}

func (p *Peer) connect(address string, isInitialConnection bool, isClient bool) {
	if contains(p.Peerlist.Peers, address) {
		return
	}
	if address == p.Ip+":"+p.Port {
		return
	}
	// Attempt to connect
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	p.Peerlist.Peers = append(p.Peerlist.Peers, address)

	// create and store encoder for connection
	encoder := gob.NewEncoder(conn)
	p.ConnectionEncoderStruct.ConnectionEncoderList = append(p.ConnectionEncoderStruct.ConnectionEncoderList, encoder)

	// create decoder for connection
	decoder := gob.NewDecoder(conn)

	if isInitialConnection {
		// Ask connection for addresses
		p.massConnect(encoder, decoder)
	}

	if isClient {
		e := encoder.Encode("JoinedNetwork")
		if e != nil {
			fmt.Println(e.Error())
		}
		e = encoder.Encode(p.Ip + ":" + p.Port)
		if e != nil {
			fmt.Println(e.Error())
		}
	}

	// Start message decoding loop for connection
	go p.handleConnection(encoder, decoder)
}

func (p *Peer) massConnect(encoder *gob.Encoder, decoder *gob.Decoder) {
	addresses := askPeerForAddresses(encoder, decoder)

	// needs to loop through the addresses and connect
	for _, address := range addresses.Peers {
		p.connect(address, false, true)
	}
}

func askPeerForAddresses(encoder *gob.Encoder, decoder *gob.Decoder) PeerList {
	// ask for GetPeerList
	e := encoder.Encode("GetAddresses")
	if e != nil {
		fmt.Println(e.Error())
	}
	// the returned list of addresses (Peerlist)
	var addresses PeerList
	e = decoder.Decode(&addresses)
	if e != nil {
		fmt.Println(e.Error())
		panic(-1)
	}
	return addresses
}

func (p *Peer) listen() (string, string) {
	// Create listener
	fmt.Println("Initializing peer")
	ln, err := net.Listen("tcp", ":")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error creating tcp listener. Terminating...")
		panic(-1)
	}
	// Get ip and port from listener
	go p.listenLoop(ln)

	return util.GetHostInfo(ln)
}

// accepts incoming connections
func (p *Peer) listenLoop(ln net.Listener) {
	defer ln.Close()
	for {
		conn, e := ln.Accept()
		if e != nil {
			fmt.Println(e.Error())
			panic(-1)
		}
		fmt.Println("connection received")

		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)

		go p.handleConnection(encoder, decoder)
	}
}

func (p *Peer) FloodMessage() {
	for _, encoderElement := range p.ConnectionEncoderStruct.ConnectionEncoderList {
		encoderElement.Encode("Hello")
		fmt.Println("Sending Hello")
	}
}

func (p *Peer) FloodTransaction(from string, to string, amount int) {
	for i, encoderElement := range p.ConnectionEncoderStruct.ConnectionEncoderList {
		encoderElement.Encode("FloodTransaction")

		// Making transaction struct
		var t account.Transaction
		t.From = from
		t.To = to
		t.Amount = amount

		// Making the transaction also on local peer and making sure it doeesn't take money multiple times
		if i == len(p.ConnectionEncoderStruct.ConnectionEncoderList)-1 {
			p.Ledger.Accounts[t.From] -= amount
			p.Ledger.Accounts[t.To] += amount
		}

		// Sending the transaction struct
		encoderElement.Encode(t)

	}
}

func (p *Peer) FloodAuthTransaction(from string, to string, amount int) {
	for i, encoderElement := range p.ConnectionEncoderStruct.ConnectionEncoderList {
		encoderElement.Encode("FloodAuthTransaction")

		// Making SignedTransaction struct
		var t account.SignedTransaction
		t.From = from
		t.To = to
		t.Amount = amount

		// Converting int to string
		amountString := strconv.Itoa(t.Amount)

		// Making all fields into one string
		allFieldsString := t.From + t.To + amountString

		// Making public key (modulus) into a big int
		modulus := new(big.Int)
		modulus, error := modulus.SetString(p.ID, 10)
		if !error {
			fmt.Println("SetString: error")
			return
		}

		// msgToBeSignede (all fields) into big.Int
		msgToBeSigned := new(big.Int)
		msgToBeSigned, error2 := msgToBeSigned.SetString(allFieldsString, 10)
		if !error2 {
			fmt.Println("SetString: error2")
			return
		}

		// Making the signature and turning into a string
		signature := encrypt.CreateSignature(msgToBeSigned, p.SecretKey, modulus)

		// put signature in the SignedTransaction struct
		t.Signature = signature.String()

		// Local transaction if p.ID (pk) is equal to t.From
		if i == len(p.ConnectionEncoderStruct.ConnectionEncoderList)-1 && p.ID == t.From {
			p.Ledger.Accounts[t.From] -= amount
			p.Ledger.Accounts[t.To] += amount
		}

		// Sending the transaction struct
		encoderElement.Encode(t)

	}
}

func (p *Peer) handleConnection(encoder *gob.Encoder, decoder *gob.Decoder) {
	for {
		// decoding the message
		var msg string
		e := decoder.Decode(&msg)
		if e != nil {
			fmt.Println(e.Error())
			panic(-1)
		}

		// switch/case
		switch msg {
		case "FloodTransaction":
			var t *account.Transaction
			e = decoder.Decode(&t)
			if e != nil {
				fmt.Println(e.Error())
				panic(-1)
			}
			p.Ledger.Transaction(t)

		case "Hello":
			fmt.Println("this is a message flood test with peer port:", p.Port)

		case "FloodAuthTransaction":
			var t *account.SignedTransaction
			e = decoder.Decode(&t)
			if e != nil {
				fmt.Println(e.Error())
				panic(-1)
			}
			p.Ledger.SignedTransaction(t)

		case "JoinedNetwork":
			var newPeer string
			e := decoder.Decode(&newPeer)
			if e != nil {
				fmt.Println(e.Error())
				panic(-1)
			}
			p.connect(newPeer, false, false)
		case "GetAddresses":
			e := encoder.Encode(&p.Peerlist)
			if e != nil {
				fmt.Println(e.Error())
			}

		default:
			fmt.Println("received unknown mesg", msg)
		}

	}
}

/*
P2P part two Backlog

	// Dealing with the many type conversions: createSignature() and verifySignature() uses big ints not strings
	- Create method for turning big.Int into strings (first into bytes and then characters) **(CHECK)
	- Create method for turning strings back into big.Int **CHECK
	- Create a single method for converting amount (int) into a string **CHECK

	// Key generation and test cases
	- import encrypt package **CHECK
	- Do keygen in initPeer() **CHECK
	- remove flood.transaction() from test cases and insert p1.floodAuthTransaction() **CHECK

	// Setting up signedTransaction
	- Adding the new transaction classes to transaction folder **CHECK

	// floodAuthTransaction() method
	- Make a new method floodAuthTransaction() which reuses parts from floodTransaction() and subsequent case which has 'FloodAuthenticatedTransaction'
		- creates signature from p.D, concatString and p.ID (only the modulus part of the public key) **CHECK
		- locally verify signature with the from-pk **(CHECK)
		- encode and remotely verify signature with the from-pk **(CHECK)


	Note: the verifySignature() is probably placed in signedTransaction class

*/
