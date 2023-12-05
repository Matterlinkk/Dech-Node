package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	"github.com/Matterlinkk/Dech-Wallet/structs"
	"time"
)

type Transaction struct {
	TransactionID string
	Time          time.Time
	From          *structs.Point
	To            *structs.Point
	Data          string
	Nonce         uint32
	Signature     structs.Signature
}

func CreateTransaction(sender, receiver *user.User, data string, signature structs.Signature) *Transaction {
	timing := time.Now()
	message := fmt.Sprintf("%v%v%v%v", timing, sender.Id, receiver.Id, sender.Nonce)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	sender.IncreaseNonce()

	return &Transaction{
		TransactionID: transactionHash,
		Time:          timing,
		From:          sender.Id,
		To:            receiver.Id,
		Data:          data,
		Signature:     signature,
		Nonce:         sender.Nonce,
	}
}

func TransactionsKeyPair(tnx []Transaction) (string, error) {
	jsonData, err := json.Marshal(tnx)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
