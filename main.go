package main

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"math/big"
)

func main() {

	alice := user.CreateUser(big.NewInt(124421125165))
	bob := user.CreateUser(big.NewInt(412412412))

	db := block.CreateBlockchain()

	sig1 := signature.SignMessage("message", alice.GetKeys())
	tnx1 := transaction.CreateTransaction(alice, bob, "message", *sig1)

	sig2 := signature.SignMessage("q", bob.GetKeys())
	tnx2 := transaction.CreateTransaction(bob, alice, "q", *sig2)

	sot := make([]transaction.Transaction, 0)
	sot = append(sot, *tnx1)
	sot = append(sot, *tnx2)

	block1 := block.CreateBlock(sot, db)
	block2 := block.CreateBlock(sot, db)
	db.AddBlock(*block1)
	db.AddBlock(*block2)

	db.KeyPair()
}
