package addressbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type AddressBook struct {
	AddressBook map[string]string //map[nickname]address
	sync.Mutex
}

func createAddressBook() *AddressBook {
	return &AddressBook{
		AddressBook: make(map[string]string),
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
		if os.IsNotExist(err) {
			return createAddressBook(), nil
		}
		log.Printf("Error ReadFile: %s", err)
		return &AddressBook{}, fmt.Errorf("error reading file: %s", err)
	}

	// Check if the file is empty
	if len(fileContent) == 0 {
		// File is empty, create a new AddressBook
		data := createAddressBook()
		if err := saveJSON(data, filename); err != nil {
			return &AddressBook{}, fmt.Errorf("error creating file: %s", err)
		}
		return data, nil
	}

	data := createAddressBook()

	err = json.Unmarshal(fileContent, &data.AddressBook)
	if err != nil {
		return &AddressBook{}, fmt.Errorf("error parsing JSON: %s", err)
	}

	return data, nil
}

func AddKeyValue(key, value, filename string) {
	data, err := LoadJSON(filename)

	if err != nil {
		log.Panicf("Error in json-loading: %s", err)
		return
	}

	log.Print("After LoadJSON: ", data)

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	data.AddressBook[key] = value

	fmt.Print("After LoadJSON: ", data)

	err = saveJSON(data, filename)
	if err != nil {
		log.Printf("saveJSON error: %s", err)
	}

	fmt.Printf("Added new user: %s=%s\n", key, value)
}
