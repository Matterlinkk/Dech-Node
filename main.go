package main

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/blockchain"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/user"
	walletOp "github.com/Matterlinkk/Dech-Wallet/operations"
	"math/big"
)

func main() {

	blockchainArray := blockchain.CreateBlockchain()

	sashkaPrivate := big.NewInt(132132)
	sashka := user.CreateUser(sashkaPrivate)

	ilyaPrivate := big.NewInt(132321231231312)
	ilya := user.CreateUser(ilyaPrivate)

	secret := walletOp.GetSharedSecret(sashka.Id, ilya.PrivateKey)

	m := "message"
	encryptedM := walletOp.GetEncryptedMessage(secret, m)

	signature1, _ := walletOp.SignMessage(encryptedM, ilya.PrivateKey, *ilya.Id)
	tnx1 := transaction.CreateTransaction(ilya, sashka, encryptedM, *signature1)

	signature2, _ := walletOp.SignMessage(encryptedM, sashka.PrivateKey, *sashka.Id)
	tnx2 := transaction.CreateTransaction(sashka, ilya, encryptedM, *signature2)

	b := block.CreateBlock()
	b.AddTnx(tnx1)

	b1 := block.CreateBlock()
	b1.AddTnx(tnx2)

	blockchainArray.AddBlock(*b1)
	blockchainArray.AddBlock(*b)

	blockchainArray.KeyPair()
}
