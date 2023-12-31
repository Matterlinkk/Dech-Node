package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/addressbook"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/hash"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/operations"
	"github.com/Matterlinkk/Dech-Wallet/publickey"
	s "github.com/Matterlinkk/Dech-Wallet/signature"
	"time"
	"unicode"
)

type Transaction struct {
	TransactionID string
	Time          time.Time
	FromAddress   string
	FromPublicKey publickey.PublicKey
	ToAddress     string
	ToPublicKey   publickey.PublicKey
	DataType      string
	Data          string
	Nonce         uint32
	Signature     s.Signature
}

func CreateTransaction(sender *user.User, receiver addressbook.TransactionReceiver, data []byte, signature s.Signature) *Transaction {

	verify := s.VerifySignature(signature, string(data), sender.PublicKey.PublicKey)

	if !verify {
		return nil
	}

	timing := time.Now()

	secret := keys.GetSharedSecret(receiver.PublicKey, sender.PrivateKey)

	encryptedMessage := operations.GetEncryptedString(secret.Bytes(), string(data))

	message := fmt.Sprintf("%v%v%v%v%s", timing, sender.Id, receiver.Id, sender.Nonce, encryptedMessage)
	transactionHash := fmt.Sprintf("%x", hash.SHA1(message))

	dt := checkType(data)
	if dt == "unknown" {
		return nil
	}

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

func isText(text []byte) bool {
	for _, char := range string(text) {
		if !unicode.IsGraphic(char) {
			return false
		}
	}
	return true
}

func checkType(data []byte) string {
	if isText(data) {
		return "text"
	} else if bytes.Equal(data[:2], []byte{255, 216}) {
		return ".jpg"
	} else if bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return ".png"
	} else if bytes.Equal(data[:3], []byte{73, 68, 51}) {
		return ".mp3"
	} else if bytes.Equal(data[:8], []byte{0, 0, 0, 24, 102, 116, 121, 112}) {
		return ".mp4"
	}
	return "unknown"
}
