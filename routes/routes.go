package routes

import (
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/handlers"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterRoutes(r chi.Router, blockchain *block.Blockchain, txChan chan transaction.Transaction, messageMap *map[string][]message.Message, loggedUser *user.User) {

	r.Post("/user/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r)
	})

	r.Get("/user/find/{user}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FindUser(w, r)
	})

	r.Get("/user/login/{nickname}", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, loggedUser)

	})

	r.Get("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowUserProfile(w, r, loggedUser)
	})

	r.Post("/tnx/create/{receiver}/media", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTnxWithMultimedia(w, r, txChan, *loggedUser)
	})

	r.Get("/tnx/create/{receiver}/text", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTnxWithText(w, r, txChan, *loggedUser)
	})

	r.Get("/blockchain/show", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowBlockchain(w, r, blockchain)
	})

	r.Get("/message/show/{from}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMessage(w, r, *messageMap, *loggedUser)
	})

	r.Get("/addressbook/show", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowAddressBook(w, r, "addressbook.json")
	})
}

func CallHandler(address string) (int, string) {
	response, err := http.Get(address)

	if err != nil {
		log.Printf("Error getting endpoint: %s", err)
		return 0, ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return response.StatusCode, ""
	}

	return response.StatusCode, string(body)
}
