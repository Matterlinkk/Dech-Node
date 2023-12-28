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
	FromAddress   string
	FromPublicKey publickey.PublicKey
	ToAddress     string
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

	secret := keys.GetSharedSecret(receiver.PublicKey, sender.GetPrivate())

	encryptedMessage := operations.GetEncryptedMessage(secret, data)

	message := fmt.Sprintf("%v%v%v%v%s", timing, sender.Id, receiver.Id, sender.ShowNonce(), encryptedMessage)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	dt := CheckType(data)

	sender.IncreaseNonce()

	return &Transaction{
		TransactionID: transactionHash,
		Time:          timing,
		FromAddress:   sender.Id,
		FromPublicKey: sender.PublicKey,
		ToAddress:     receiver.Id,
		ToPublicKey:   receiver.PublicKey,
		Data:          encryptedMessage,
		DataType:      dt,
		Signature:     signature,
		Nonce:         sender.ShowNonce(),
	}
}

func TransactionsKeyPair(tnx []Transaction) (string, error) {
	jsonData, err := json.Marshal(tnx)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func CheckType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	default:
		return "Unknown type"
	}
	//need finalize
}
