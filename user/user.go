package user

import (
	walletOperations "github.com/Matterlinkk/Dech-Wallet/operations"
	"github.com/Matterlinkk/Dech-Wallet/structs"
	"math/big"
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

func (u *User) CreateUser(privateKey *big.Int) *User {
	keyPair := walletOperations.GetKeyPair(privateKey)

	return &User{
		Id:         keyPair.PublicKey,
		Nonce:      0,
		PublicKey:  keyPair.PublicKey,
		PrivateKey: keyPair.PrivateKey,
		Balance:    0,
	}
}
