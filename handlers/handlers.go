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

func AddTnx(w http.ResponseWriter, r *http.Request, tnxChannel chan transaction.Transaction, sender user.User) {
	receiverStr := chi.URLParam(r, "receiver")

	addressBook, err := addressbook.LoadJSON("addressbook.json")
	if err != nil {
		log.Printf("Error json-loading: %s", err)
		return
	}

	receiver, ok := addressBook.AddressBook[receiverStr]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		responseText := fmt.Sprintf("User %s does not exist", receiverStr)
		w.Write([]byte(responseText))
		return
	}

	message := r.URL.Query().Get("data")

	signature := signature.SignMessage(message, sender.GetKeys())

	tnx := transaction.CreateTransaction(&sender, receiver, message, *signature)

	transportchan.TnxToBlock(tnxChannel, *tnx)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Transaction successfully sent"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	pkString := r.URL.Query().Get("pk")
	nickname := r.URL.Query().Get("nickname")
	password := r.URL.Query().Get("password")
	pK, _ := new(big.Int).SetString(pkString, 10)

	newUser := user.CreateUser(pK, nickname)

	filename := fmt.Sprintf("%s.txt", nickname)

	user.CreateUserFile(filename, password, *newUser)

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
}

func GetMessage(w http.ResponseWriter, r *http.Request, messageMap map[string][]message.Message) {
	fromStr := chi.URLParam(r, "from")
	toStr := chi.URLParam(r, "to")

	addressBook, err := addressbook.LoadJSON("addressbook.json")
	if err != nil {
		log.Printf("Error json-loading: %s", err)
		return
	}

	from, ok := addressBook.AddressBook[fromStr]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}
	to, ok := addressBook.AddressBook[toStr]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}

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

	entries := make(map[string]addressbook.TransactionReceiver)

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

func LoginUser(w http.ResponseWriter, r *http.Request, loggedUser *user.User) {
	nickname := chi.URLParam(r, "nickname")
	password := r.URL.Query().Get("password")

	filename := fmt.Sprintf("%s.txt", nickname)

	*loggedUser = user.ReadUserFile(filename, password, w)

	if loggedUser.Nickname == "" {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	inputText := fmt.Sprintf("Access accepted user is: %s", loggedUser.Nickname)
	w.Write([]byte(inputText))
}

func ShowUserProfile(w http.ResponseWriter, r *http.Request, loggedUser *user.User) {
	jsonResponse, err := json.Marshal(loggedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating JSON response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}
