package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/operations"
	"github.com/Matterlinkk/Dech-Wallet/publickey"
	s "github.com/Matterlinkk/Dech-Wallet/signature"
	"time"
)

type Transaction struct {
	TransactionID string
	Time          time.Time
	FromAdress    string
	FromPublicKey publickey.PublicKey
	ToAdress      string
	ToPublicKey   publickey.PublicKey
	Data          string
	DataType      string
	Nonce         uint32
	Signature     s.Signature
}

func CreateTransaction(sender, receiver *user.User, data string, signature s.Signature) *Transaction {

	verify := s.VerifySignature(signature, data, sender.PublicKey.PublicKey)

	if !verify {
		return nil
	}

	timing := time.Now()

	secret := keys.GetSharedSecret(receiver.PublicKey, sender.PrivateKey)

	encryptedMessage := operations.GetEncryptedMessage(secret, data)

	message := fmt.Sprintf("%v%v%v%v%s", timing, sender.Id, receiver.Id, sender.Nonce, encryptedMessage)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	dt := checkType(data)

	sender.IncreaseNonce()

	return &Transaction{
		TransactionID: transactionHash,
		Time:          timing,
		FromAdress:    sender.Id,
		FromPublicKey: sender.PublicKey,
		ToAdress:      receiver.Id,
		ToPublicKey:   receiver.PublicKey,
		Data:          encryptedMessage,
		DataType:      dt,
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

func checkType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	default:
		return "Unknown type"
	}
	//need finalize
}
