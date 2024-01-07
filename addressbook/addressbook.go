package addressbook

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/publickey"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type TransactionReceiver struct {
	Id        string              `json:"id"`
	PublicKey publickey.PublicKey `json:"publicKey"`
}

type AddressBook struct {
	AddressBook map[string]TransactionReceiver // map[nickname]struct{Id string; PublicKey publickey.PublicKey}
	sync.Mutex
}

func createAddressBook() *AddressBook {
	return &AddressBook{
		AddressBook: make(map[string]TransactionReceiver),
	}
}

func saveJSON(data *AddressBook, filename string) error {
	jsonData, err := json.Marshal(data.AddressBook)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}

func LoadJSON(filename string) (*AddressBook, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		//Check if data exists
		if os.IsNotExist(err) {
			return createAddressBook(), nil
		}
		log.Printf("Error ReadFile: %s", err)
		return &AddressBook{}, fmt.Errorf("error reading data: %s", err)
	}

	// Check if the data is empty
	if len(fileContent) == 0 {

		// File is empty, create a new AddressBook
		data := createAddressBook()
		if err := saveJSON(data, filename); err != nil {
			return &AddressBook{}, fmt.Errorf("Error creating data: %s", err)
		}
		return data, nil
	}

	data := createAddressBook()

	err = json.Unmarshal(fileContent, &data.AddressBook)
	if err != nil {
		return &AddressBook{}, fmt.Errorf("Error parsing JSON: %s", err)
	}

	return data, nil
}

func AddKeyValue(nickname string, user user.User, filename string) {
	data, err := LoadJSON(filename)

	if err != nil {
		log.Panicf("Error in json-loading: %s", err)
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	data.AddressBook[nickname] = struct {
		Id        string              `json:"id"`
		PublicKey publickey.PublicKey `json:"publicKey"`
	}{
		Id:        user.Id,
		PublicKey: user.PublicKey,
	}

	err = saveJSON(data, filename)
	if err != nil {
		log.Printf("saveJSON error: %s", err)
	}

	fmt.Printf("Added new user: %s\n", nickname)
}
