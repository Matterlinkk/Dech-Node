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
	Nonce      uint32
	Nickname   string
	PublicKey  publickey.PublicKey
	PrivateKey *big.Int
	Balance    uint32
}

func CreateUser(privateKey *big.Int, nickname string) *User {
	keyPair := keys.GetKeys(privateKey)

	return &User{
		Id:         keyPair.PublicKey.GetAdress(),
		Nickname:   nickname,
		Nonce:      0,
		PublicKey:  *keyPair.PublicKey,
		PrivateKey: keyPair.PrivateKey,
		Balance:    0,
	}
}

func (u *User) GetKeys() keys.KeyPair {
	return keys.KeyPair{
		PrivateKey: u.PrivateKey,
		PublicKey:  &u.PublicKey,
	}
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
	u.Nonce++
	return u.Nonce
}

func FindByUser(nickname string, userDB []User) User {
	for _, u := range userDB {
		if u.Nickname == nickname {
			return u
		}
	}
	return User{}
}
