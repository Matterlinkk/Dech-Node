package message

import (
	"github.com/Matterlinkk/Dech-Node/addressbook"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/operations"
	"strings"
	"time"
)

type Message struct {
	dataType string
	data     []byte
	date     time.Time
}

func CreateMessage(data []byte, transaction transaction.Transaction) Message {

	return Message{
		dataType: transaction.DataType,
		data:     data,
		date:     transaction.Time,
	}
}

func (message Message) ShowData(user user.User, fromUser addressbook.TransactionReceiver) string {

	secret := keys.GetSharedSecret(fromUser.PublicKey, user.PrivateKey)

	return operations.GetDecryptedString(secret.Bytes(), string(message.data))
}

func (message Message) ShowString(user user.User, fromUser addressbook.TransactionReceiver) string {

	text := message.date.String()

	index := strings.Index(text, ".")

	secret := keys.GetSharedSecret(fromUser.PublicKey, user.PrivateKey)

	return text[:index] + " " + operations.GetDecryptedString(secret.Bytes(), string(message.data))
}

func (message Message) ShowDataType() string {
	return message.dataType
}

func CreateMessageMap() *map[string][]Message {
	messageMap := make(map[string][]Message)
	return &messageMap
}

func ParseBlock(messageMap *map[string][]Message, block block.Block) *map[string][]Message {

	for _, tx := range block.SetOfTransactions {

		message := CreateMessage([]byte(tx.Data), tx)
		key := tx.FromAddress + ":" + tx.ToAddress
		(*messageMap)[key] = append((*messageMap)[key], message)
	}
	return messageMap
}
