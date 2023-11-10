package structs

import "time"

type User struct {
	Id        string
	Nonce     uint32
	PublicKey string
	SecretKey string
	Balance   uint32
}

func (u *User) IncreaseNonce() uint32 {
	u.Nonce++
	return u.Nonce
}

type Transaction struct {
	TransactionID string
	Time          time.Time
	From          string
	To            string
	Nonce         uint32
}

type Block struct {
	BlockId           string
	PrevHash          string
	SetOfTransactions []Transaction
}
