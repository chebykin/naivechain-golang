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
	blockchain   []*Block
)

type Block struct {
	Index        int64
	PreviousHash string
	Timestamp    int
	Data         []byte
	Hash         string
}

func getGenesisBlock() *Block {
	return &Block{
		Index:        0,
		PreviousHash: "0",
		Timestamp:    1498381610,
		Data:         []byte("my genesis block!!!"),
		Hash:         "816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
	}
}

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
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", http_port),
		Handler: r,
	}

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