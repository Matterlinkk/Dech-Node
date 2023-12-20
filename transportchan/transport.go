package transportchan

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"time"
)

func BlockToBlockchain(SetOfBlock chan block.Block, block block.Block) {
	SetOfBlock <- block
}

func TnxToBlock(setOfTnx chan transaction.Transaction, tnx transaction.Transaction) {
	setOfTnx <- tnx
}

func ProcessBlock(chanBlock chan block.Block, chanTnx chan transaction.Transaction, db *block.Blockchain, messageMap *map[string][]message.Message) {
	for {
		set := block.CreateSetOfTransactions()

		var timer *time.Timer

	CollectTransactions:
		for {
			if timer == nil {
				timer = time.NewTimer(5 * time.Second)
			} else {
				timer.Reset(5 * time.Second)
			}

			select {
			case tnx := <-chanTnx:
				set = append(set, tnx)
			case <-timer.C:
				newBlock := block.CreateBlock(set, db)
				fmt.Println("Block: ", newBlock)
				messageMap = message.ParseBlock(messageMap, *newBlock)
				fmt.Println("MessageMap: ", messageMap)
				BlockToBlockchain(chanBlock, *newBlock)
				break CollectTransactions
			}
		}
	}
}

func ProcessBlockchain(setOfBlock chan block.Block, blockchain *block.Blockchain) {
	for b := range setOfBlock {
		blockchain.Mutex.Lock()
		blockchain.AddBlock(b)
		blockchain.Unlock()
	}
}
