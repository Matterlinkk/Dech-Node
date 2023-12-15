package handlers

import (
	"encoding/json"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/Matterlinkk/Dech-Wallet/signature"
	"github.com/go-chi/chi"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func AddTnx(w http.ResponseWriter, r *http.Request, userDB []user.User, tnxChannel chan transaction.Transaction) {
	receiverIdStr := chi.URLParam(r, "receiver") //http://localhost:8080/tnx/create/1/0/message?data=wqe
	receiverId, _ := strconv.Atoi(receiverIdStr)
	senderIdStr := chi.URLParam(r, "sender")
	senderId, _ := strconv.Atoi(senderIdStr)
	message := r.URL.Query().Get("data")

	signature := signature.SignMessage(message, userDB[senderId].GetKeys())

	tnx := transaction.CreateTransaction(&userDB[senderId], &userDB[receiverId], message, *signature)

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
	pK, _ := new(big.Int).SetString(pkString, 10)

	newUser := user.CreateUser(pK)

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
