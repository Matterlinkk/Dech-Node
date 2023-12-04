package structs

import (
	"github.com/Matterlinkk/Dech-Wallet/structs"
	"math/big"
	"time"
)

type User struct {
	Id         *structs.Point
	Nonce      uint32
	PublicKey  *structs.Point
	PrivateKey *big.Int
	Balance    uint32
}

func (u *User) IncreaseNonce() uint32 {
	u.Nonce++
	return u.Nonce
}

type Transaction struct {
	TransactionID string
	Time          time.Time
	From          *structs.Point
	To            *structs.Point
	Data          string
	Nonce         uint32
	Signature     structs.Signature
}

type Block struct {
	BlockId           string
	PrevHash          string
	SetOfTransactions []Transaction
}
