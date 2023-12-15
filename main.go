package main

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/handlers"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	var UserDB []user.User

	db := block.CreateBlockchain()

	channelTnx := make(chan transaction.Transaction)
	channelBlock := make(chan block.Block)

	go func() {
		transportchan.ProcessBlockchain(channelBlock, db)
	}()

	go func() {
		transportchan.ProcessBlock(channelBlock, channelTnx, db)
	}()

	r.Get("/user/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r, &UserDB)
	})
	r.Get("/user/list", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowUserDatabase(w, r, UserDB)
	})
	r.Get("/tnx/create/{sender}/{receiver}/message", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTnx(w, r, UserDB, channelTnx)
	})
	r.Get("/blockchain/show", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowBlockchain(w, r, db)
	})

	http.ListenAndServe(":8080", r)
}

//http://localhost:8080/user/create?pk=1
//http://localhost:8080/user/create?pk=12321442112545512512
//http://localhost:8080/tnx/create/1/0/message?data=wqe
