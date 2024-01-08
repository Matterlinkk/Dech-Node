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
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// @Summary Creates a transaction with a text message
// @Description Creates a transaction with a text message for transmission from sender to receiver.
// @ID AddTnxWithText
// @Tags transactions
// @Param receiver path string true "Receiver's name"
// @Param data query string true "Text data"
// @Success 201 {string} string "Transaction successfully sent"
// @Failure 400 {string} string "Request error"
// @Failure 409 {string} string "Transaction declined"
// @Router /tnx/create/{receiver}/text [get]
func AddTnxWithText(w http.ResponseWriter, r *http.Request, tnxChannel chan transaction.Transaction, sender user.User) {
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

	tnx := transaction.CreateTransaction(&sender, receiver, []byte(message), *signature)
	if tnx == nil {
		w.WriteHeader(http.StatusConflict)
		responseText := fmt.Sprintf("Transaction declined", receiverStr)
		w.Write([]byte(responseText))
		return
	}

	transportchan.TnxToBlock(tnxChannel, *tnx)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Transaction successfully sent"))
}

// UploadMedia is a handler that serves an HTML form for uploading media.
// @Summary Serves an HTML form for uploading media.
// @Description This endpoint serves an HTML form for users to upload media files.
// @ID UploadMedia
// @Tags media
// @Success 200 {string} string "HTML form successfully served"
// @Failure 500 {string} string "Internal Server Error"
// @Router /upload/media [get]
func UploadMedia(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir("Dech-Node"))
	if err != nil {
		http.Error(w, "Unable to determine file path", http.StatusInternalServerError)
		return
	}

	htmlPath := filepath.Join(dir, "html", "upload-file.html")

	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read HTML file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	w.Write(htmlContent)
}

// @Summary Creates a transaction with a media file
// @Description Creates a transaction with a media file for transmission from sender to receiver.
// @ID AddTnxWithMultimedia
// @Tags transactions
// @Param receiver path string true "Recipient's name"
// @Param file formData file true "Multimedia file for transmission"
// @Success 201 {string} string "Transaction successfully sent"
// @Failure 400 {string} string "Request error"
// @Failure 404 {string} string "User not found"
// @Failure 413 {string} string "The file is too big"
// @Failure 409 {string} string "Transaction declined(unknown file format)"
// @Failure 500 {string} string "Internal server error"
// @Router /tnx/create/{receiver}/media [post]
func AddTnxWithMultimedia(w http.ResponseWriter, r *http.Request, tnxChannel chan transaction.Transaction, sender user.User) {

	const maxUploadSize = 15 << 20

	receiverStr := r.FormValue("receiver")

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

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize) //15MB max size

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		errorStr := fmt.Sprintf("ParseMultipartForm error: %s", err)
		http.Error(w, errorStr, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error when retrieving a file from a form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Println("Filename: ", handler.Filename, "\nFile size: ", handler.Size, "\nReceiver's name: ", receiverStr)

	if handler.Size > maxUploadSize {
		http.Error(w, "The file is too big", http.StatusBadRequest)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseText := "Error reading file data"
		w.Write([]byte(responseText))
		return
	}

	fmt.Println("Bytes: ", fileBytes[:40])

	signature := signature.SignMessage(string(fileBytes), sender.GetKeys())

	tnx := transaction.CreateTransaction(&sender, receiver, fileBytes, *signature)
	if tnx == nil {
		w.WriteHeader(http.StatusConflict)
		responseText := fmt.Sprintf("Transaction declined(unknown file format)", receiverStr)
		w.Write([]byte(responseText))
		return
	}

	transportchan.TnxToBlock(tnxChannel, *tnx)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Transaction successfully sent"))
}

// CreateUser creates a new user with the provided public key, nickname, and password.
// @Summary Create a new user
// @Description Endpoint to create a new user with the provided public key, nickname, and password.
// @ID CreateUser
// @Tags users
// @Param pk query string true "Public key of the user (base-10 string)"
// @Param nickname query string true "Nickname for the new user"
// @Param password query string true "Password for the new user"
// @Success 200 {string} string "User successfully created"
// @Failure 400 {string} string "Bad Request"
// @Router /user/create [post]
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

// ShowBlockchain returns the blockchain data in a plain text format.
// @Summary Show blockchain data
// @Description Endpoint to retrieve and display the blockchain data in plain text format.
// @ID ShowBlockchain
// @Tags blockchain
// @Success 200 {string} string "Blockchain data successfully retrieved"
// @Router /blockchain/show [get]
func ShowBlockchain(w http.ResponseWriter, r *http.Request, db *block.Blockchain) {

	dbJson, _ := json.Marshal(db)

	dbString := strings.ReplaceAll(string(dbJson), ",", "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(dbString))
}

// FindUser retrieves information about a user based on the provided username.
// @Summary Find user by username
// @Description Endpoint to find and display user information based on the provided username.
// @ID FindUser
// @Tags users
// @Param user path string true "Username of the user to find"
// @Success 200 {string} string "User information successfully retrieved"
// @Failure 404 {string} string "User not found"
// @Router /user/find/{user} [get]
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

// GetMessage retrieves messages between the logged-in user and the specified sender.
// @Summary Get messages between users
// @Description Endpoint to retrieve and display messages between the logged-in user and the specified sender.
// @ID GetMessage
// @Tags messages
// @Param from path string true "Username of the sender"
// @Success 200 {string} string "Messages successfully retrieved"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error encoding JSON"
// @Router /message/show/{from} [get]
func GetMessage(w http.ResponseWriter, r *http.Request, messageMap map[string][]message.Message, loggedUser user.User) {
	fromStr := chi.URLParam(r, "from")

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

	key := from.Id + ":" + loggedUser.Id
	messages, ok := messageMap[key]
	if !ok {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}

	var messagesStrings []string

	for _, msg := range messages {
		messagesStrings = append(messagesStrings, msg.ShowString(loggedUser, from))
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

// ShowAddressBook returns the address book entries in JSON format.
// @Summary Get address book entries
// @Description Endpoint to retrieve and display entries from the address book in JSON format.
// @ID ShowAddressBook
// @Tags address book
// @Success 200 {string} string "Address book entries successfully retrieved"
// @Failure 500 {string} string "Error creating JSON response"
// @Router /addressbook/show [get]
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

// LoginUser authenticates a user based on the provided nickname and password.
// @Summary Authenticate user
// @Description Endpoint to authenticate a user based on the provided nickname and password.
// @ID LoginUser
// @Tags authentication
// @Param nickname path string true "User nickname"
// @Param password query string true "User password"
// @Success 200 {string} string "Access accepted"
// @Failure 401 {string} string "Access denied"
// @Router /user/login/{nickname} [get]
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

// ShowUserProfile retrieves and displays the profile of the logged-in user.
// @Summary Show user profile
// @Description Endpoint to retrieve and display the profile of the logged-in user.
// @ID ShowUserProfile
// @Tags users
// @Success 200 {string} string "User profile retrieved successfully"
// @Failure 401 {string} string "User not authenticated"
// @Router /user/profile [get]
func ShowUserProfile(w http.ResponseWriter, r *http.Request, loggedUser *user.User) {
	jsonResponse, err := json.Marshal(loggedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating JSON response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}

func GetLastMediaFile(w http.ResponseWriter, r *http.Request, loggedUser user.User, messageMap map[string][]message.Message) {
	fromStr := chi.URLParam(r, "from")

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

	key := from.Id + ":" + loggedUser.Id
	messages, ok := messageMap[key]
	if !ok {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}

	var message message.Message

	for _, msg := range messages {
		if msg.ShowDataType() != "text" {
			message = msg
		}
	}

	bytes := []byte(message.ShowData(loggedUser, from))

	contentType := http.DetectContentType(bytes)

	w.Header().Set("Content-Disposition", "attachment; filename=")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
