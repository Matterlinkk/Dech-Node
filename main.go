package main

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/handlers"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	messageMap := message.CreateMessageMap()

	var UserDB []user.User

	db := block.CreateBlockchain()

	channelTnx := make(chan transaction.Transaction)
	channelBlock := make(chan block.Block)

	go func() {
		transportchan.ProcessBlockchain(channelBlock, db)
	}()

	go func() {
		transportchan.ProcessBlock(channelBlock, channelTnx, db, messageMap)
	}()

	r.Get("/user/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r, &UserDB)
	})
	r.Get("/user/list", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowUserDatabase(w, r, UserDB)
	})

	r.Get("/user/find/{user}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FindUser(w, r, UserDB)
	})

	r.Get("/tnx/create/{sender}/{receiver}/message", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTnx(w, r, UserDB, channelTnx)
	})

	r.Get("/blockchain/show", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowBlockchain(w, r, db)
	})

	r.Get("/message/show/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMessage(w, r, UserDB, *messageMap)
	})

	http.ListenAndServe(":8080", r)
}

//http://localhost:8080/user/create?pk=123&nickname=Alice
//http://localhost:8080/user/create?pk=321&nickname=Bob
//http://localhost:8080/tnx/create/Alice/Bob/message?data=some_message
//http://localhost:8080/blockchain/show
