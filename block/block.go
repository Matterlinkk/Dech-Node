package block

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"log"
)

type Block struct {
	BlockId           string
	PrevHash          string
	SetOfTransactions []transaction.Transaction
}

func CreateBlock() *Block {
	return &Block{
		BlockId:           "",
		PrevHash:          "",
		SetOfTransactions: make([]transaction.Transaction, 0),
	}
}

func (block *Block) AddTnx(tnx *transaction.Transaction) {
	verify := signature.VerifySignature(tnx.Signature, tnx.Data, tnx.FromPublicKey.PublicKey)
	if verify {
		block.SetOfTransactions = append(block.SetOfTransactions, *tnx)
	}
}

func (block *Block) AddBlock(blockchain *[]Block) {

	jsonString, err := transaction.TransactionsKeyPair(block.SetOfTransactions)

	if err != nil {
		log.Panicf("Ошибка при преобразовании в JSON: %v", err)
	}

	if len(*blockchain) == 0 {
		blockId := hash.SHA1("" + jsonString)
		block.BlockId = fmt.Sprintf("%x", blockId)
		block.PrevHash = ""
		*blockchain = append(*blockchain, *block)
		return
	}

	// Calculate the new block's ID and set the PrevHash
	blockId := hash.SHA1((*blockchain)[len(*blockchain)-1].BlockId + jsonString)
	block.BlockId = fmt.Sprintf("%x", blockId)
	block.PrevHash = (*blockchain)[len(*blockchain)-1].BlockId

	// Append the new block to the blockchain
	*blockchain = append(*blockchain, *block)
}

func (block *Block) KeyPair() string {
	jsonData, _ := json.Marshal(block)
	return string(jsonData)
}

func BlockKeyPair(block []Block) {
	for i := range block {
		jsonData, _ := json.Marshal(block[i])
		fmt.Println(string(jsonData) + "\n")
	}
}
