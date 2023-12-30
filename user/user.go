package user

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/operations"
	"github.com/Matterlinkk/Dech-Wallet/publickey"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

type User struct {
	Id         string
	Nonce      uint32
	Nickname   string
	PublicKey  publickey.PublicKey
	PrivateKey *big.Int
}

func CreateUser(privateKey *big.Int, nickname string) *User {
	keyPair := keys.GetKeys(privateKey)

	return &User{
		Id:         keyPair.PublicKey.GetAdress(),
		Nickname:   nickname,
		Nonce:      0,
		PublicKey:  *keyPair.PublicKey,
		PrivateKey: keyPair.PrivateKey,
	}
}

func CreateUserFile(filename, password string, user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	encryptedUser := operations.GetEncryptedString([]byte(password), string(jsonData))

	filePath := fmt.Sprintf("user/userfiles/%s", filename)

	err = ioutil.WriteFile(filePath, []byte(encryptedUser), 0644)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}

func ReadUserFile(filename, password string, w http.ResponseWriter) User {
	filePath := fmt.Sprintf("user/userfiles/%s", filename)

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Invalid nickname", http.StatusNotFound)
		return User{}
	}

	decryptedData := operations.GetDecryptedString([]byte(password), string(fileData))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decrypting data: %s", err), http.StatusNotFound)
		return User{}
	}

	var user User
	err = json.Unmarshal([]byte(decryptedData), &user)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusNotFound)
		return User{}
	}

	return user
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
