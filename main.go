package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	http_port    string
	p2p_port     string
	initialPeers []string
)

func initHttpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		// TODO: wrap with locks or use channels instead
		bc, err := json.Marshal(blockchain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bc)
	}).Methods("GET")

	r.HandleFunc("/mineBlock", func(w http.ResponseWriter, r *http.Request) {
		newBlock := generateNextBlock([]byte(r.FormValue("data")))

		addBlock(newBlock)

		// TODO: broadcast

		fmt.Println("block added >>>", newBlock)

		w.WriteHeader(200)
	}).Methods("POST")

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", http_port),
		Handler: r,
	}

	fmt.Printf("Listening on port http://localhost:%s\n", http_port)
	s.ListenAndServe()
}

func main() {
	http_port = os.Getenv("HTTP_PORT")
	if http_port == "" {
		http_port = "3001"
	}

	p2p_port = os.Getenv("P2P_PORT")
	if p2p_port == "" {
		p2p_port = "6001"
	}

	blockchain = []*Block{getGenesisBlock()}

	initHttpServer()

	fmt.Println(blockchain)
	fmt.Println("vim-go")
}
