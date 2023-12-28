package message

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"strings"
	"time"
)

type Message struct {
	data string
	date time.Time
}

func CreateMessage(data string, transaction transaction.Transaction) Message {
	return Message{
		data: data,
		date: transaction.Time,
	}
}

func (message Message) ShowString() string {

	text := message.date.String()

	index := strings.Index(text, ".")

	return text[:index] + " " + message.data
}

func CreateMessageMap() *map[string][]Message {
	messageMap := make(map[string][]Message)
	return &messageMap
}

func ParseBlock(messageMap *map[string][]Message, block block.Block) *map[string][]Message {

	for _, tx := range block.SetOfTransactions {

		message := CreateMessage(tx.Data, tx)
		key := tx.FromAdress + ":" + tx.ToAdress
		(*messageMap)[key] = append((*messageMap)[key], message)
	}
	return messageMap
}