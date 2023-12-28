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

func RegisterRoutes(r chi.Router, UserDB []user.User, blockchain *block.Blockchain, txChan chan transaction.Transaction, messageMap *map[string][]message.Message) {

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
		handlers.AddTnx(w, r, UserDB, txChan)
	})

	r.Get("/blockchain/show", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowBlockchain(w, r, blockchain)
	})

	r.Get("/message/show/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMessage(w, r, UserDB, *messageMap)
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
