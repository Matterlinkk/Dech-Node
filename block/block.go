package block

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	"log"
	"sync"
)

type Block struct {
	BlockId           string
	PrevHash          string
	SetOfTransactions []transaction.Transaction
}

type Blockchain struct {
	sync.Mutex
	BlockArray []Block
}

func CreateBlock(sot []transaction.Transaction, bc *Blockchain) *Block {
	jsonString, err := transaction.TransactionsKeyPair(sot)

	if err != nil {
		log.Panicf("Ошибка при преобразовании в JSON: %s", err)
	}

	if len(bc.BlockArray) == 0 {
		blockId := hash.SHA1("" + jsonString)
		blockIdStr := fmt.Sprintf("%x", blockId)
		prevHash := ""
		return &Block{
			BlockId:           blockIdStr,
			PrevHash:          prevHash,
			SetOfTransactions: sot,
		}
	}

	// Calculate the new block's ID and set the PrevHash
	blockId := hash.SHA1(bc.GetLastBlock().BlockId + jsonString)
	blockIdStr := fmt.Sprintf("%x", blockId)
	prevHash := bc.GetLastBlock().BlockId

	// Append the new block to the blockchain
	return &Block{
		BlockId:           blockIdStr,
		PrevHash:          prevHash,
		SetOfTransactions: sot,
	}
}

func (block *Block) KeyPair() string {
	jsonData, _ := json.Marshal(block)
	return string(jsonData)
}

func CreateBlockchain() *Blockchain {
	return &Blockchain{
		BlockArray: make([]Block, 0),
	}
}

func (bc *Blockchain) GetLastBlock() Block {
	if len(bc.BlockArray) == 0 {
		panic("Length is 0")
	}
	return bc.BlockArray[len(bc.BlockArray)-1]
}

func (bc *Blockchain) AddBlock(block Block) {

	bc.BlockArray = append(bc.BlockArray, block)
	return
}

func (bc *Blockchain) KeyPair() {
	for _, block := range bc.BlockArray {
		jsonData, err := json.Marshal(block)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Println(string(jsonData) + "\n")
	}
}

func CreateSetOfTransactions() []transaction.Transaction {
	set := make([]transaction.Transaction, 0)
	return set
}
