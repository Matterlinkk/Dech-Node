package operations

import (
	"encoding/json"
	"fmt"
	nodeStructs "github.com/Matterlinkk/Dech-Node/structs"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	walletOperations "github.com/Matterlinkk/Dech-Wallet/operations"
	walletStructs "github.com/Matterlinkk/Dech-Wallet/structs"
	"log"
	"time"
)

func CreateTransaction(sender, receiver *nodeStructs.User, data string, signature walletStructs.Signature) *nodeStructs.Transaction {
	timing := time.Now()
	message := fmt.Sprintf("%v%v%v%v", timing, sender.Id, receiver.Id, sender.Nonce)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	sender.IncreaseNonce()

	return &nodeStructs.Transaction{
		TransactionID: transactionHash,
		Time:          timing,
		From:          sender.Id,
		To:            receiver.Id,
		Data:          data,
		Signature:     signature,
		Nonce:         sender.Nonce,
	}
}

func TransactionsKeyPair(transactions []nodeStructs.Transaction) (string, error) {
	jsonData, err := json.Marshal(transactions)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func CreateBlock() *nodeStructs.Block {
	return &nodeStructs.Block{
		BlockId:           "",
		PrevHash:          "",
		SetOfTransactions: make([]nodeStructs.Transaction, 0),
	}
}

func AddTnx(block *nodeStructs.Block, tnx *nodeStructs.Transaction) {
	verify := walletOperations.VerifySignature(&tnx.Signature, tnx.Data, tnx.From)
	if verify {
		block.SetOfTransactions = append(block.SetOfTransactions, *tnx)
	}
}

func AddBlock(blockchain *[]nodeStructs.Block, block nodeStructs.Block) {

	jsonString, err := TransactionsKeyPair(block.SetOfTransactions)

	if err != nil {
		log.Panicf("Ошибка при преобразовании в JSON: %v", err)
	}

	if len(*blockchain) == 0 {
		blockId := hash.SHA1("" + jsonString)
		block.BlockId = fmt.Sprintf("%x", blockId)
		block.PrevHash = ""
		*blockchain = append(*blockchain, block)
		return
	}

	// Calculate the new block's ID and set the PrevHash
	blockId := hash.SHA1((*blockchain)[len(*blockchain)-1].BlockId + jsonString)
	block.BlockId = fmt.Sprintf("%x", blockId)
	block.PrevHash = (*blockchain)[len(*blockchain)-1].BlockId

	// Append the new block to the blockchain
	*blockchain = append(*blockchain, block)
}

func BlockKeyPair(block []nodeStructs.Block) {
	for i := range block {
		jsonData, _ := json.Marshal(block[i])
		fmt.Println(string(jsonData) + "\n")
	}
}
