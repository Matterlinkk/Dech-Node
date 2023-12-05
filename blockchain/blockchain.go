package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Wallet/hash"
)

type Blockchain struct {
	BlockArray []block.Block
}

func CreateBlockchain() *Blockchain {
	return &Blockchain{
		BlockArray: make([]block.Block, 0),
	}
}

func (bc *Blockchain) GetLastBlock() block.Block {
	if len(bc.BlockArray) == 0 {
		panic("Length is 0")
	}
	return bc.BlockArray[len(bc.BlockArray)-1]
}

func (bc *Blockchain) AddBlock(block block.Block) {

	if len(bc.BlockArray) == 0 {
		hashValue := hash.SHA1(block.KeyPair())
		block.BlockId = fmt.Sprintf("%x", string(hashValue[:]))
		block.PrevHash = ""
		bc.BlockArray = append(bc.BlockArray, block)
		return
	}

	hashValue := hash.SHA1(bc.GetLastBlock().BlockId + block.KeyPair())
	block.BlockId = fmt.Sprintf("%x", string(hashValue[:]))
	block.PrevHash = bc.GetLastBlock().BlockId
	bc.BlockArray = append(bc.BlockArray, block)
	return

}

func (bc *Blockchain) KeyPair() {
	for _, blk := range bc.BlockArray {
		jsonData, _ := json.Marshal(blk)
		fmt.Println(string(jsonData))
		fmt.Println()
	}
}
