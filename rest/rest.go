package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	blockchain "github.com/devholic77/duckcoin/blockchanin"
	"github.com/devholic77/duckcoin/utils"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u urlDescription) String() string {
	return "URL Description"
}

type addBlockBody struct {
	Message string `'json:"message"`
}

type errResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         "/blocks",
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         "/blocks/{hash}",
			Method:      "GET",
			Description: "GET A Block",
		},
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain())
	case "POST":
		var block addBlockBody
		fmt.Println()
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&block))
		blockchain.BlockChain().AddBlock(block.Message)
		fmt.Println(block)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		hash := vars["hash"]
		block := blockchain.BlockChain().GetBlock(hash)
		json.NewEncoder(rw).Encode(block)
	}

}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")

	fmt.Printf("RESTAPI server Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
