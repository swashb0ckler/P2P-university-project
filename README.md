# P2P Golang 

A **Peer to Peer** golang implementation running a ledger using RSA signatures for authenticated transactions.  In order to use it run the TestLedger(), ignore the console logs and instead see if the accounts have the right amount after being flooded. 

<p align="center">
<img width="300px" src="https://images.unsplash.com/photo-1639322537228-f710d846310a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1632&q=80"/>
</p>


###How it Works
A peer can join the network and connect to the others peers already available. When a peer is started it connects to another peer (if there is one) and gets their peer list whereafter they mass connect to everyone in the network. Each peer is thus on entry into the network connecting to all other peers and these are connecting back. The peers are using Go routines to listen concurrently for new peers on the network as well as waiting for new transactions. 

Transactions are flooded so every peer will have the same version of the ledger and idea of which accounts have what stored. Furthermore, the transactions are signed, so only the account which is authorized to make a transaction can make a transaction. 

###Try It out 

Run the **go test** and the if the ledgers are identical and if the accounts have the right amount after making signed transactions 

```
	// Every ledger is identical
	assert.Equal(T, p1.Ledger.Accounts, p2.Ledger.Accounts, p3.Ledger.Accounts, p4.Ledger.Accounts, "The two peers ledgers should be the same.")

	// p1 should have: (-200) + (-200) == -400
	assert.Equal(T, p1.Ledger.Accounts[p1.ID], -400, "p1 should have the right balance")

	// p2 should have: 100 + 300 + 300 + == 700
	assert.Equal(T, p1.Ledger.Accounts[p2.ID], 700, "p2 should have the right balance")

	// p3 should have: (-100) + (-300) + (-300) + == -700
	assert.Equal(T, p1.Ledger.Accounts[p3.ID], -700, "p3 should have the right balance")

	// p4 should have: 200 + 200  == 400
	assert.Equal(T, p1.Ledger.Accounts[p4.ID], 400, "p4 should have the right balance")```

```