package operations

import (
	"encoding/json"
	"fmt"
	"log"
	"project/hash"
	"project/structs"
	"time"
)

func CreateTransaction(sender, receiver *structs.User) *structs.Transaction {
	timing := time.Now()
	message := fmt.Sprintf("%v%v%v%v", timing, sender.Id, receiver.Id, sender.Nonce)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	sender.IncreaseNonce()

	return &structs.Transaction{
		TransactionID: transactionHash,
		Time:          timing,
		From:          sender.Id,
		To:            receiver.Id,
		Nonce:         sender.Nonce,
	}
}

func TransactionsKeyPair(transactions []structs.Transaction) (string, error) {
	jsonData, err := json.Marshal(transactions)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func CreateBlock(blockchain []structs.Block, seroftnx []structs.Transaction) *structs.Block {

	jsonString, err := TransactionsKeyPair(seroftnx)

	if err != nil {
		log.Panicf("Ошибка при преобразовании в JSON:", err)
	}

	if len(blockchain) == 0 {

		blockId := hash.SHA1("" + jsonString)

		return &structs.Block{
			BlockId:           fmt.Sprintf("%x", blockId),
			PrevHash:          "",
			SetOfTransactions: seroftnx,
		}
	}

	blockId := hash.SHA1(blockchain[len(blockchain)-1].BlockId + jsonString)

	return &structs.Block{
		BlockId:           fmt.Sprintf("%x", blockId),
		PrevHash:          blockchain[len(blockchain)-1].BlockId,
		SetOfTransactions: seroftnx,
	}
}

func BlockKeyPair(block []structs.Block) {
	for i := range block {
		jsonData, _ := json.Marshal(block[i])
		fmt.Println(string(jsonData) + "\n")
	}
}
