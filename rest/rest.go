package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	blockchain "github.com/devholic77/duckcoin/blockchain"
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

type balanceResponse struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type addTxPayLoad struct {
	To     string
	Amount int
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
			URL:         "/status",
			Method:      "GET",
			Description: "See status of blockchain",
		},
		{
			URL:         "/blocks",
			Method:      "GET",
			Description: "Get All Block",
		},
		{
			URL:         "/blocks",
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "message:string",
		},
		{
			URL:         "/blocks/{hash}",
			Method:      "GET",
			Description: "GET A Block",
		},
		{
			URL:         "/balance/{address}",
			Method:      "GET",
			Description: "GET TxOuts for an Address",
		},
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		hash := vars["hash"]
		block, err := blockchain.FindBlock(hash)
		utils.HandleErr(err)
		json.NewEncoder(rw).Encode(block)
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")

	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		response := balanceResponse{
			Address: address,
			Amount:  amount,
		}
		json.NewEncoder(rw).Encode(response)
	default:
		json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain()))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool))
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayLoad
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errResponse{"not enough funds"})
	}
	rw.WriteHeader(http.StatusCreated)
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
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")

	fmt.Printf("RESTAPI server Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
