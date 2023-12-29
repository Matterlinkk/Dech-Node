package user

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/publickey"
	"math/big"
)

type User struct {
	Id         string
	nonce      uint32
	Nickname   string
	PublicKey  publickey.PublicKey
	privateKey *big.Int
}

func CreateUser(privateKey *big.Int, nickname string) *User {
	keyPair := keys.GetKeys(privateKey)

	return &User{
		Id:         keyPair.PublicKey.GetAdress(),
		Nickname:   nickname,
		nonce:      0,
		PublicKey:  *keyPair.PublicKey,
		privateKey: keyPair.PrivateKey,
	}
}

func (u *User) GetKeys() keys.KeyPair {
	return keys.KeyPair{
		PrivateKey: u.privateKey,
		PublicKey:  &u.PublicKey,
	}
}

func (u *User) GetPrivate() *big.Int {
	return u.privateKey
}

func (u *User) KeyPair() {
	jsonData, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonData) + "\n")
}

func (u *User) IncreaseNonce() uint32 {
	u.nonce++
	return u.nonce
}

func (u *User) ShowNonce() uint32 {
	return u.nonce
}
