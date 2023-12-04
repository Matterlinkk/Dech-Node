package main

import (
	nodeOp "github.com/Matterlinkk/Dech-Node/operations"
	nodeStruct "github.com/Matterlinkk/Dech-Node/structs"
	walletOp "github.com/Matterlinkk/Dech-Wallet/operations"
	"math/big"
)

func main() {

	blockchain := make([]nodeStruct.Block, 0)

	sashkaPrivate := big.NewInt(132132)
	sashka := nodeOp.CreateUser(sashkaPrivate)

	ilyaKeys := big.NewInt(13213132)
	ilya := nodeOp.CreateUser(ilyaKeys)

	secret := walletOp.GetSharedSecret(sashka.Id, ilya.PrivateKey)

	m := "message"
	encryptedM := walletOp.GetEncryptedMessage(secret, m)

	signature1, _ := walletOp.SignMessage(encryptedM, ilya.PrivateKey, *ilya.Id)
	tnx1 := nodeOp.CreateTransaction(ilya, sashka, encryptedM, *signature1)

	signature2, _ := walletOp.SignMessage(encryptedM, sashka.PrivateKey, *sashka.Id)
	tnx2 := nodeOp.CreateTransaction(sashka, ilya, encryptedM, *signature2)

	block := nodeOp.CreateBlock()
	nodeOp.AddTnx(block, tnx1)

	block1 := nodeOp.CreateBlock()
	nodeOp.AddTnx(block1, tnx2)

	nodeOp.AddBlock(&blockchain, *block)
	nodeOp.AddBlock(&blockchain, *block1)

	nodeOp.BlockKeyPair(blockchain)
}
