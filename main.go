package main

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"math/big"
	"time"
)

func main() {

	alice := user.CreateUser(big.NewInt(124421125165))
	bob := user.CreateUser(big.NewInt(412412412))

	db := block.CreateBlockchain()

	sig1 := signature.SignMessage("first", alice.GetKeys())
	tnx1 := transaction.CreateTransaction(alice, bob, "first", *sig1)

	sig2 := signature.SignMessage("first", bob.GetKeys())
	tnx2 := transaction.CreateTransaction(bob, alice, "second", *sig2)

	sig3 := signature.SignMessage("third", bob.GetKeys())
	tnx3 := transaction.CreateTransaction(bob, alice, "third", *sig3)

	channelTnx := make(chan transaction.Transaction)
	channelBlock := make(chan block.Block)

	go func() {
		transportchan.ProcessBlockchain(channelBlock, db)
	}()

	go func() {
		transportchan.ProcessBlock(channelBlock, channelTnx, db)
	}()

	go func() {
		time.Sleep(6 * time.Second)
		transportchan.TnxToBlock(channelTnx, *tnx2)
	}()

	go func() {
		time.Sleep(21 * time.Second)
		transportchan.TnxToBlock(channelTnx, *tnx3)
	}()

	go func() {
		transportchan.TnxToBlock(channelTnx, *tnx1)
	}()

	time.Sleep(30 * time.Second)
	fmt.Println("Db: \n")
	db.KeyPair()

}
