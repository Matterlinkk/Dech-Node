package handlers

import (
	"encoding/json"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"github.com/go-chi/chi/v5"
	"math/big"
	"net/http"
	"strings"
)

func AddTnx(w http.ResponseWriter, r *http.Request, userDB []user.User, tnxChannel chan transaction.Transaction) {
	receiverStr := chi.URLParam(r, "receiver") //http://localhost:8080/tnx/create/alice/bob/message?data=wqe
	senderStr := chi.URLParam(r, "sender")
	sender := user.FindByUser(senderStr, userDB)
	receiver := user.FindByUser(receiverStr, userDB)
	message := r.URL.Query().Get("data")

	signature := signature.SignMessage(message, sender.GetKeys())

	tnx := transaction.CreateTransaction(&sender, &receiver, message, *signature)

	transportchan.TnxToBlock(tnxChannel, *tnx)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Transaction successfully sent"))
}

func ShowUserDatabase(w http.ResponseWriter, r *http.Request, userDB []user.User) {
	signatureJson, _ := json.Marshal(userDB)

	signatureString := strings.ReplaceAll(string(signatureJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(signatureString))
}

func CreateUser(w http.ResponseWriter, r *http.Request, userDB *[]user.User) {

	pkString := r.URL.Query().Get("pk") //http://localhost:8080/user/create?pk=2
	nickname := r.URL.Query().Get("nickname")
	pK, _ := new(big.Int).SetString(pkString, 10)

	newUser := user.CreateUser(pK, nickname)

	*userDB = append(*userDB, *newUser)

	w.WriteHeader(200)
	w.Write([]byte("User successfully сreated"))
}

func ShowBlockchain(w http.ResponseWriter, r *http.Request, db *block.Blockchain) {

	dbJson, _ := json.Marshal(db)

	dbString := strings.ReplaceAll(string(dbJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(dbString))
}

func FindUser(w http.ResponseWriter, r *http.Request, userDB []user.User) {

	user := chi.URLParam(r, "user")

	var result string

	for _, userStruct := range userDB {
		if user == userStruct.Nickname {
			userJson, _ := json.Marshal(userStruct)
			result = strings.ReplaceAll(string(userJson), ",", "\n")
		}
	}

	if len(result) != 0 {
		w.WriteHeader(200)
		w.Write([]byte(result))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
	}
}

func GetMessage(w http.ResponseWriter, r *http.Request, userDB []user.User, messageMap map[string][]message.Message) {
	fromStr := chi.URLParam(r, "from")
	toStr := chi.URLParam(r, "to")

	from := user.FindByUser(fromStr, userDB)
	to := user.FindByUser(toStr, userDB)

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
