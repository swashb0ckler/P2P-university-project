package test

import (
	"fmt"
	"p2p/peer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLedger(T *testing.T) {
	// 1. Create some Peers
	fmt.Println("------printed from peer 1------")
	p1, ip1, port1 := peer.InitPeer("")
	fmt.Println("peer:", p1)
	fmt.Println("ip:", ip1)
	fmt.Println("port:", port1)
	fmt.Println("---------------------")

	time.Sleep(1 * time.Second)
	fmt.Println("------printed from peer 2 -------")
	p2, ip2, port2 := peer.InitPeer(ip1 + ":" + port1)
	fmt.Println("peer:", p2)
	fmt.Println("ip:", ip2)
	fmt.Println("port:", port2)
	fmt.Println("---------------------")

	time.Sleep(1 * time.Second)
	fmt.Println("------printed from peer 3 -------")
	p3, ip3, port3 := peer.InitPeer(ip2 + ":" + port2)
	fmt.Println("peer:", p3)
	fmt.Println("ip:", ip3)
	fmt.Println("port:", port3)
	fmt.Println("---------------------")

	time.Sleep(1 * time.Second)
	fmt.Println("------printed from peer 4 -------")
	p4, ip4, port4 := peer.InitPeer(ip3 + ":" + port3)
	fmt.Println("peer:", p4)
	fmt.Println("ip:", ip4)
	fmt.Println("port:", port4)
	fmt.Println("---------------------")

	time.Sleep(2 * time.Second)
	fmt.Println("this is p1 Peerlist:", p1.Peerlist)
	fmt.Println("this is p2 Peerlist:", p2.Peerlist)
	fmt.Println("this is p3 Peerlist:", p3.Peerlist)
	fmt.Println("this is p4 Peerlist:", p4.Peerlist)

	// 2. Execute some transactions
	time.Sleep(5 * time.Second)

	// all valid transactions
	p3.FloodAuthTransaction(p3.ID, p2.ID, 100)
	p3.FloodAuthTransaction(p3.ID, p2.ID, 300)
	p3.FloodAuthTransaction(p3.ID, p2.ID, 300)
	p1.FloodAuthTransaction(p1.ID, p4.ID, 200)
	p1.FloodAuthTransaction(p1.ID, p4.ID, 200)

	// invalid transaction should not be counted towards final account balances
	p2.FloodAuthTransaction(p1.ID, p4.ID, 200)
	p3.FloodAuthTransaction(p1.ID, p4.ID, 200)
	p4.FloodAuthTransaction(p1.ID, p4.ID, 200)
	p3.FloodAuthTransaction(p2.ID, p4.ID, 200)
	p1.FloodAuthTransaction(p2.ID, p4.ID, 200)
	p1.FloodAuthTransaction(p4.ID, p4.ID, 200)

	time.Sleep(5 * time.Second)

	// 3. Assert that the ledger of every Peer is updated correctly
	fmt.Println("account 1 ------------------", p1.Ledger.Accounts)
	fmt.Println("account 2 ------------------", p2.Ledger.Accounts)
	fmt.Println("account 3 ------------------", p3.Ledger.Accounts)
	fmt.Println("account 4 ------------------", p4.Ledger.Accounts)

	time.Sleep(5 * time.Second)

	// Every ledger is identical
	assert.Equal(T, p1.Ledger.Accounts, p2.Ledger.Accounts, p3.Ledger.Accounts, p4.Ledger.Accounts, "The two peers ledgers should be the same.")

	// p1 should have: (-200) + (-200) == -400
	assert.Equal(T, p1.Ledger.Accounts[p1.ID], -400, "p1 should have the right balance")

	// p2 should have: 100 + 300 + 300 + == 700
	assert.Equal(T, p1.Ledger.Accounts[p2.ID], 700, "p2 should have the right balance")

	// p3 should have: (-100) + (-300) + (-300) + == -700
	assert.Equal(T, p1.Ledger.Accounts[p3.ID], -700, "p3 should have the right balance")

	// p4 should have: 200 + 200  == 400
	assert.Equal(T, p1.Ledger.Accounts[p4.ID], 400, "p4 should have the right balance")

}
