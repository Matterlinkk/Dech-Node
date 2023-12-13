package main

import (
	"github.com/Matterlinkk/Dech-Node/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Get("/message/sign", handlers.SignMessageHandler)
	r.Get("/user/create", handlers.CreateUser)
	r.Get("/user/list", handlers.ShowDatabase)
	r.Get("/tnx/create/{sender}/{receiver}/message", handlers.AddTnx)

	http.ListenAndServe(":8080", r)
}
