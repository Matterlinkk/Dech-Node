package main

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
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
		ProcessBlockchain(channelBlock, db)
	}()

	go func() {
		ProcessBlock(channelBlock, channelTnx, db)
	}()

	go func() {
		time.Sleep(6 * time.Second)
		TnxToBlock(channelTnx, *tnx2)
	}()

	go func() {
		time.Sleep(21 * time.Second)
		TnxToBlock(channelTnx, *tnx3)
	}()

	go func() {
		TnxToBlock(channelTnx, *tnx1)
	}()

	time.Sleep(30 * time.Second)
	fmt.Println("Db: \n")
	db.KeyPair()

}

func BlockToBlockchain(SetOfBlock chan block.Block, block block.Block) {
	SetOfBlock <- block
}

func TnxToBlock(setOfTnx chan transaction.Transaction, tnx transaction.Transaction) {
	setOfTnx <- tnx
}

func ProcessBlock(chanBlock chan block.Block, chanTnx chan transaction.Transaction, db *block.Blockchain) {
	for {
		set := block.CreateSetOfTransactions()

		var timer *time.Timer

	CollectTransactions:
		for {
			if timer == nil {
				timer = time.NewTimer(5 * time.Second)
			} else {
				timer.Reset(5 * time.Second)
			}

			select {
			case tnx := <-chanTnx:
				set = append(set, tnx)
			case <-timer.C:
				fmt.Println("Creating block")
				newBlock := block.CreateBlock(set, db)
				fmt.Println("Block: ", newBlock)
				BlockToBlockchain(chanBlock, *newBlock)
				break CollectTransactions
			}
		}
	}
}

func ProcessBlockchain(setOfBlock chan block.Block, blockchain *block.Blockchain) {
	for b := range setOfBlock {
		blockchain.Mutex.Lock()
		blockchain.AddBlock(b)
		blockchain.Unlock()
	}
}
