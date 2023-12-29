package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Matterlinkk/Dech-Node/addressbook"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"github.com/go-chi/chi/v5"
	"log"
	"math/big"
	"net/http"
	"strings"
)

func AddTnx(w http.ResponseWriter, r *http.Request, tnxChannel chan transaction.Transaction) { //fix after login logic
	receiverStr := chi.URLParam(r, "receiver")
	senderStr := chi.URLParam(r, "sender")

	addressBook, err := addressbook.LoadJSON("addressbook.json")
	if err != nil {
		log.Printf("Error json-loading: %s", err)
		return
	}

	sender := addressBook.AddressBook[senderStr]
	receiver := addressBook.AddressBook[receiverStr]

	log.Printf("sender: %v\n", sender)
	log.Printf("receiver: %v\n", receiver)

	message := r.URL.Query().Get("data")

	log.Printf("mesage: %v\n", message)

	signature := signature.SignMessage(message, sender.GetKeys())
	log.Printf("sender.GetKeys(): %v\n", sender.GetKeys())

	log.Printf("signature: %v\n", signature)

	tnx := transaction.CreateTransaction(&sender, &receiver, message, *signature)

	transportchan.TnxToBlock(tnxChannel, *tnx)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Transaction successfully sent"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	pkString := r.URL.Query().Get("pk")
	nickname := r.URL.Query().Get("nickname")
	pK, _ := new(big.Int).SetString(pkString, 10)

	newUser := user.CreateUser(pK, nickname)

	addressbook.AddKeyValue(newUser.Nickname, *newUser, "addressbook.json")

	w.WriteHeader(200)
	w.Write([]byte("User successfully сreated"))
}

func ShowBlockchain(w http.ResponseWriter, r *http.Request, db *block.Blockchain) {

	dbJson, _ := json.Marshal(db)

	dbString := strings.ReplaceAll(string(dbJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(dbString))
}

func FindUser(w http.ResponseWriter, r *http.Request) {

	user := chi.URLParam(r, "user")

	addressBook, err := addressbook.LoadJSON("addressbook.json")
	if err != nil {
		log.Printf("Error json-loading: %s", err)
		return
	}

	userStruct, ok := addressBook.AddressBook[user]

	if ok {
		userStr, _ := json.Marshal(userStruct)
		result := strings.ReplaceAll(string(userStr), ",", "\n")
		w.WriteHeader(200)
		w.Write([]byte(result))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
	}
} // refactor

func GetMessage(w http.ResponseWriter, r *http.Request, messageMap map[string][]message.Message) {
	fromStr := chi.URLParam(r, "from")
	toStr := chi.URLParam(r, "to")

	addressBook, err := addressbook.LoadJSON("addressbook.json")
	if err != nil {
		log.Printf("Error json-loading: %s", err)
		return
	}

	from := addressBook.AddressBook[fromStr]
	to := addressBook.AddressBook[toStr]

	key := from.Id + ":" + to.Id
	messages, ok := messageMap[key]
	if !ok {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}

	var messagesStrings []string

	for _, msg := range messages {
		messagesStrings = append(messagesStrings, msg.ShowString())
	}

	msgJson, err := json.Marshal(messagesStrings)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msgJson)
}

func ShowAddressBook(w http.ResponseWriter, r *http.Request, filename string) {
	data, err := addressbook.LoadJSON(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading address book: %s", err), http.StatusInternalServerError)
		return
	}

	entries := make(map[string]user.User)

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	for key, value := range data.AddressBook {
		entries[key] = value
	}

	jsonResponse, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating JSON response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}
