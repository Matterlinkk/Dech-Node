package handlers

import (
	"encoding/json"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/keys"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"github.com/go-chi/chi"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

var tnxSlice []transaction.Transaction

var db []user.User

func AddTnx(w http.ResponseWriter, r *http.Request) {
	receiverIdStr := chi.URLParam(r, "receiver") //http://localhost:8080/tnx/create/1/0/message?data=wqe
	receiverId, _ := strconv.Atoi(receiverIdStr)
	senderIdStr := chi.URLParam(r, "sender")
	senderId, _ := strconv.Atoi(senderIdStr)
	message := r.URL.Query().Get("data")

	signature := signature.SignMessage(message, db[senderId].GetKeys())

	tnx := transaction.CreateTransaction(&db[senderId], &db[receiverId], message, *signature)

	transaction, _ := json.Marshal(tnx)

	transactionString := strings.ReplaceAll(string(transaction), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(transactionString))
}

func ShowDatabase(w http.ResponseWriter, r *http.Request) {
	signatureJson, _ := json.Marshal(db)

	signatureString := strings.ReplaceAll(string(signatureJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(signatureString))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	pkString := r.URL.Query().Get("pk") //http://localhost:8080/user/create?pk=2
	pK, _ := new(big.Int).SetString(pkString, 10)

	idStr := r.URL.Query().Get("id")

	newUser := user.CreateUser(pK, idStr)

	db = append(db, *newUser)

	w.WriteHeader(200)
	w.Write([]byte("Success"))
}

func SignMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message") //http://localhost:8080/message/sign/?message=q&pk=2

	privateKey := r.URL.Query().Get("pk")

	pK, _ := new(big.Int).SetString(privateKey, 10)

	key := keys.GetKeys(pK)

	signature := signature.SignMessage(message, key)

	log.Print(signature.GetSignature)

	signatureJson, _ := json.Marshal(signature)

	signatureString := strings.ReplaceAll(string(signatureJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(signatureString))

}
